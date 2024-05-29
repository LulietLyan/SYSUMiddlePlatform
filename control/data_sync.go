package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func getAllProjectTable() (configList []*models.ProjectTable, err error) {
	if e := mysql.DB.Find(&configList).Error; e != nil {
		return nil, e
	}
	return configList, nil
}

func getSourceConfig(config *models.ProjectTable, str1 string) (str string) {
	str = fmt.Sprintf("CREATE TABLE flink_%s.flink_source_%s "+
		"%s"+ //主键必须标明NOT ENFORCED
		" WITH (\n"+
		"  'connector'  = 'mysql-cdc',\n"+
		"  'hostname'   = '%s',\n"+
		"  'port'   = '%d',\n"+
		"  'database-name'   = '%s',\n"+
		"  'table-name' = '%s',\n"+
		"  'username'   = '%s',\n"+
		"  'password'   = '%s'\n"+
		")",
		config.PT_remote_db_name, config.PT_remote_table_name, str1,
		config.PT_remote_hostname, config.PT_remote_port, config.PT_remote_db_name,
		config.PT_remote_table_name, config.PT_remote_username, config.PT_remote_password)
	return str
}

func getTargetConfig(config *models.ProjectTable, str1 string) (str string) {
	str = fmt.Sprintf("CREATE TABLE flink_%s.flink_target_%s "+
		"%s"+ //主键必须标明NOT ENFORCED
		" WITH (\n"+
		"  'connector'  = 'jdbc',\n"+
		"  'driver'     = 'com.mysql.cj.jdbc.Driver',\n"+
		"  'url'        = 'jdbc:mysql://47.121.29.57:3307/flink_target',\n"+
		"  'table-name' = '%d_%s_%s',\n"+
		"  'username'   = 'root',\n"+
		"  'password'   = '654321'\n"+
		")",
		config.PT_remote_db_name, config.PT_remote_table_name, str1,
		config.PU_uid, config.PT_remote_db_name, config.PT_remote_table_name)
	return str
}

func InitDataSync() error {
	todoList, err := getAllProjectTable()
	if err != nil {
		return err
	}
	// 运行所有配置信息
	for index, _ := range todoList {
		// 数据同步——获取表结构
		tmp, err := gorm.Open("mysql",
			fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				todoList[index].PT_remote_username, todoList[index].PT_remote_password,
				todoList[index].PT_remote_hostname, todoList[index].PT_remote_port, todoList[index].PT_remote_db_name))
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
			return err
		}
		QueryColumnSQL := fmt.Sprintf("SELECT "+
			"COLUMN_NAME, DATA_TYPE, COLUMN_DEFAULT, IS_NULLABLE, COLUMN_KEY, COLUMN_COMMENT "+
			"FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '%s' "+
			"AND TABLE_NAME = '%s' ORDER BY ORDINAL_POSITION;",
			todoList[index].PT_remote_db_name,
			todoList[index].PT_remote_table_name)
		var columns []struct {
			ColumnName    string `gorm:"column:COLUMN_NAME"`
			DataType      string `gorm:"column:DATA_TYPE"`
			ColumnDefault string `gorm:"column:COLUMN_DEFAULT"`
			IsNullable    string `gorm:"column:IS_NULLABLE"`
			ColumnKey     string `gorm:"column:COLUMN_KEY"`
			ColumnComment string `gorm:"column:COLUMN_COMMENT"`
		}
		if err := tmp.Raw(QueryColumnSQL).Scan(&columns).Error; err != nil {
			log.Fatalf("failed to query columns: %v", err)
			return err
		}
		for _, column := range columns {
			fmt.Println(column)
		}
		QueryCreateTableSQL := fmt.Sprintf("SHOW CREATE TABLE %s", todoList[index].PT_remote_table_name)
		var createSQLResult []struct {
			TableName   string `gorm:"column:Table"`
			CreateTable string `gorm:"column:Create Table"`
		}
		if err := tmp.Raw(QueryCreateTableSQL).Scan(&createSQLResult).Error; err != nil {
			fmt.Printf("Error executing SQL: %v\n", err)
			return err
		}
		leftIndex := strings.Index(createSQLResult[0].CreateTable, "(")
		rightIndex := strings.LastIndex(createSQLResult[0].CreateTable, ")")
		configStr := createSQLResult[0].CreateTable[leftIndex : rightIndex+1]
		fmt.Println(configStr)
		re := regexp.MustCompile(`PRIMARY KEY \((.*?)\)`)
		configStr1 := re.ReplaceAllStringFunc(configStr, func(s string) string {
			return s + " NOT ENFORCED"
		})
		re = regexp.MustCompile(`\s*DEFAULT\s+[^,]*\s*,`)
		configStr1 = re.ReplaceAllString(configStr1, ",")
		configStr1 = strings.ReplaceAll(configStr1, "datetime", "timestamp")
		fmt.Println(configStr1)
		if err := tmp.Close(); err != nil {
			return err
		}
		// 数据同步——创建数据表结构
		dropTableSQL := fmt.Sprintf("DROP TABLE IF EXISTS %d_%s_%s;",
			todoList[index].PU_uid, todoList[index].PT_remote_db_name, todoList[index].PT_remote_table_name)
		if err := mysql.DB_flink.Exec(dropTableSQL).Error; err != nil {
			log.Fatalf("failed to drop table if exists: %v", err)
			return err
		}
		createTableSQL := fmt.Sprintf("CREATE TABLE %d_%s_%s %s;",
			todoList[index].PU_uid, todoList[index].PT_remote_db_name, todoList[index].PT_remote_table_name, configStr)
		if err := mysql.DB_flink.Exec(createTableSQL).Error; err != nil {
			log.Fatalf("failed to create table: %v", err)
			return err
		}
		// 数据同步——java程序执行
		strCreatDatabase := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS flink_%s",
			todoList[index].PT_remote_db_name)
		strSource := getSourceConfig(todoList[index], configStr1)
		strTarget := getTargetConfig(todoList[index], configStr1)
		strInsert := fmt.Sprintf("INSERT INTO flink_%s.flink_target_%s select * from flink_%s.flink_source_%s\n",
			todoList[index].PT_remote_db_name, todoList[index].PT_remote_table_name,
			todoList[index].PT_remote_db_name, todoList[index].PT_remote_table_name)
		fmt.Println(strCreatDatabase)
		fmt.Println(strSource)
		fmt.Println(strTarget)
		fmt.Println(strInsert)
		// 运行java命令
		cmd := exec.Command("java",
			"-cp", "flinkdemo.jar;flink_libs/*",
			"com.demo.flink.FlinkCdcMySql",
			strCreatDatabase,
			strSource,
			strTarget,
			strInsert,
		)
		if err = cmd.Start(); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func NewProjectTable(c *gin.Context) {
	type table struct {
		ProjectId    uint   `json:"project_id"`
		DatabaseName string `json:"database_name"`
		TableName    string `json:"table_name"`
		TableConfig  string `json:"table_config"`
		HostName     string `json:"host_name"`
		UserName     string `json:"user_name"`
		Password     string `json:"password"`
		Port         uint   `json:"port"`
	}
	var t table

	if e := c.ShouldBindJSON(&t); e == nil {
		var projectTable models.ProjectTable
		projectTable = models.ProjectTable{
			PU_uid:               t.ProjectId,
			PT_remote_db_name:    t.DatabaseName,
			PT_remote_table_name: t.TableName,
			PT_remote_hostname:   t.HostName,
			PT_remote_username:   t.UserName,
			PT_remote_password:   t.Password,
			PT_remote_port:       t.Port,
		}

		tx := mysql.DB.Begin()
		if e := tx.Create(&projectTable).Error; e != nil {
			tx.Rollback()
			response.Fail(c, nil, "插入新项目表时出错")
			return
		}
		if e := tx.Commit().Error; e != nil {
			response.Fail(c, nil, "提交事务时出错")
			return
		}
		response.Success(c, gin.H{"project_id": t.ProjectId}, "")
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
