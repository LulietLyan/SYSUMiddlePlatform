package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var processMap sync.Map

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

func getEsConfig(config *models.ProjectTable, str1 string) (str string) {
	str = fmt.Sprintf("CREATE TABLE flink_%s.flink_target_%s "+
		"%s"+ //主键必须标明NOT ENFORCED
		" WITH (\n"+
		"  'connector'  = 'elasticsearch-7',\n"+
		"  'hosts'     = 'http://47.120.73.205:9200',\n"+
		"  'index'        = '%d_%s_%s_index'\n"+
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

		re := regexp.MustCompile(`PRIMARY KEY \((.*?)\)`)
		configStr1 := re.ReplaceAllStringFunc(configStr, func(s string) string {
			return s + " NOT ENFORCED"
		})
		Index1 := strings.Index(configStr1, "NOT ENFORCED")
		// 没有主键
		if Index1 == -1 {
			Index2 := strings.Index(configStr1, "KEY")
			if Index2 != -1 {
				configStr = configStr1[0:Index2-4] + "\n)"
			}

			lines := strings.Split(configStr, "\n")
			for i, line := range lines {
				// 使用Fields函数获取单词列表
				words := strings.Fields(line)
				// 只保留前两个词
				if len(words) >= 2 {
					lines[i] = "  " + words[0] + " " + words[1]
					if i < len(lines)-2 {
						lines[i] = lines[i] + ","
					}
				} else if len(words) == 1 {
					lines[i] = words[0]
				} else {
					lines[i] = ""
				}
			}
			configStr1 = strings.Join(lines, "\n")
		} else {
			Index2 := strings.Index(configStr1, "NOT ENFORCED")
			configStr = configStr1[0:Index2] + "\n)"
			configStr1 = configStr1[0:Index2+12] + "\n)"

			lines := strings.Split(configStr1, "\n")
			for i, line := range lines {
				// 使用Fields函数获取单词列表
				words := strings.Fields(line)
				// 只保留前两个词
				if len(words) >= 2 {
					if i == len(lines)-2 {
						continue
					}
					lines[i] = "  " + words[0] + " " + words[1] + ","
				} else if len(words) == 1 {
					lines[i] = words[0]
				} else {
					lines[i] = ""
				}
			}
			configStr1 = strings.Join(lines, "\n")
		}
		configStr1 = strings.ReplaceAll(configStr1, "datetime", "timestamp")

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
		// 运行java命令——同步到mysql
		cmd := exec.Command("java",
			"-cp", "mysql2mysql.jar;flink_libs/*",
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
		// 运行java命令——同步到es
		cmd1 := exec.Command("java",
			"-cp", "mysql2es.jar;flink_libs/*",
			"com.demo.flink.FlinkCdcMySql",
			strCreatDatabase,
			strSource,
			getEsConfig(todoList[index], configStr1),
			strInsert,
		)
		if err = cmd1.Start(); err != nil {
			log.Fatal(err)
			return err
		}
		// 保存Process以便后续控制
		process := cmd.Process
		process1 := cmd1.Process
		processMap.Store(todoList[index].PT_uid, [2]*os.Process{process, process1})
	}
	return nil
}

func NewProjectTable(c *gin.Context) {
	type table struct {
		Name            string `json:"name"`
		Description     string `json:"description"`
		Uid             uint   `json:"uid"`
		RemoteDbName    string `json:"remote_db_name"`
		RemoteTableName string `json:"remote_table_name"`
		RemoteHostname  string `json:"remote_hostname"`
		RemoteUsername  string `json:"remote_username"`
		RemotePassword  string `json:"remote_password"`
		RemotePort      uint   `json:"remote_port"`
	}
	var t table

	if e := c.ShouldBindJSON(&t); e == nil {
		var pu_uid uint
		if data, ok := c.Get("pu_uid"); !ok {
			response.Fail(c, nil, "没有从token解析出所需信息")
			return
		} else {
			pu_uid = data.(uint)
		}
		t.Uid = pu_uid
		// 判断是否已存在
		var projectTable1 models.ProjectTable
		if e := mysql.DB.Where("PU_uid = ? and PT_remote_db_name = ? and PT_remote_table_name = ?",
			t.Uid, t.RemoteDbName, t.RemoteTableName).First(&projectTable1).Error; e == nil {
			response.Fail(c, nil, "该项目表已存在!")
			return
		}
		// 新增记录
		var projectTable models.ProjectTable
		projectTable = models.ProjectTable{
			PT_name:              t.Name,
			PT_description:       t.Description,
			PU_uid:               t.Uid,
			PT_remote_db_name:    t.RemoteDbName,
			PT_remote_table_name: t.RemoteTableName,
			PT_remote_hostname:   t.RemoteHostname,
			PT_remote_username:   t.RemoteUsername,
			PT_remote_password:   t.RemotePassword,
			PT_remote_port:       t.RemotePort,
		}
		tx := mysql.DB.Begin()
		if e := tx.Create(&projectTable).Error; e != nil {
			tx.Rollback()
			response.Fail(c, nil, "插入新项目表时出错")
			return
		}
		if e := tx.Commit().Error; e != nil {
			response.Fail(c, nil, "插入新项目表时出错")
			return
		}
		// 数据同步——获取表结构
		tmp, err := gorm.Open("mysql",
			fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				t.RemoteUsername, t.RemotePassword, t.RemoteHostname, t.RemotePort, t.RemoteDbName))
		if err != nil {
			response.Fail(c, nil, "数据同步出错1")
			return
		}
		QueryCreateTableSQL := fmt.Sprintf("SHOW CREATE TABLE %s", t.RemoteTableName)
		var createSQLResult []struct {
			TableName   string `gorm:"column:Table"`
			CreateTable string `gorm:"column:Create Table"`
		}
		if err := tmp.Raw(QueryCreateTableSQL).Scan(&createSQLResult).Error; err != nil {
			response.Fail(c, nil, "数据同步出错2")
			return
		}
		leftIndex := strings.Index(createSQLResult[0].CreateTable, "(")
		rightIndex := strings.LastIndex(createSQLResult[0].CreateTable, ")")
		configStr := createSQLResult[0].CreateTable[leftIndex : rightIndex+1]

		re := regexp.MustCompile(`PRIMARY KEY \((.*?)\)`)
		configStr1 := re.ReplaceAllStringFunc(configStr, func(s string) string {
			return s + " NOT ENFORCED"
		})
		re = regexp.MustCompile(`\s*DEFAULT\s+[^,]*\s*,`)
		configStr1 = re.ReplaceAllString(configStr1, ",")
		configStr1 = strings.ReplaceAll(configStr1, "datetime", "timestamp")

		var projectTable2 models.ProjectTable
		if e := mysql.DB.Where("PU_uid = ? and PT_remote_db_name = ? and PT_remote_table_name = ?",
			t.Uid, t.RemoteDbName, t.RemoteTableName).First(&projectTable2).Error; e != nil {
			response.Fail(c, nil, "数据同步出错3")
			return
		}
		if err := tmp.Close(); err != nil {
			response.Fail(c, nil, "数据同步出错4")
			return
		}
		// 数据同步——创建数据表结构
		dropTableSQL := fmt.Sprintf("DROP TABLE IF EXISTS %d_%s_%s;", t.Uid, t.RemoteDbName, t.RemoteTableName)
		if err := mysql.DB_flink.Exec(dropTableSQL).Error; err != nil {
			response.Fail(c, nil, "数据同步出错5")
			return
		}
		createTableSQL := fmt.Sprintf("CREATE TABLE %d_%s_%s %s;", t.Uid, t.RemoteDbName, t.RemoteTableName, configStr)
		if err := mysql.DB_flink.Exec(createTableSQL).Error; err != nil {
			response.Fail(c, nil, "数据同步出错6")
			return
		}
		// 数据同步——java程序执行
		strCreatDatabase := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS flink_%s", t.RemoteDbName)
		strSource := fmt.Sprintf("CREATE TABLE flink_%s.flink_source_%s "+
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
			projectTable2.PT_remote_db_name, projectTable2.PT_remote_table_name, configStr1,
			projectTable2.PT_remote_hostname, projectTable2.PT_remote_port, projectTable2.PT_remote_db_name,
			projectTable2.PT_remote_table_name, projectTable2.PT_remote_username, projectTable2.PT_remote_password)
		strTarget := fmt.Sprintf("CREATE TABLE flink_%s.flink_target_%s "+
			"%s"+ //主键必须标明NOT ENFORCED
			" WITH (\n"+
			"  'connector'  = 'jdbc',\n"+
			"  'driver'     = 'com.mysql.cj.jdbc.Driver',\n"+
			"  'url'        = 'jdbc:mysql://47.121.29.57:3307/flink_target',\n"+
			"  'table-name' = '%d_%s_%s',\n"+
			"  'username'   = 'root',\n"+
			"  'password'   = '654321'\n"+
			")",
			projectTable2.PT_remote_db_name, projectTable2.PT_remote_table_name, configStr1,
			projectTable2.PU_uid, projectTable2.PT_remote_db_name, projectTable2.PT_remote_table_name)
		strInsert := fmt.Sprintf("INSERT INTO flink_%s.flink_target_%s select * from flink_%s.flink_source_%s\n",
			t.RemoteDbName, t.RemoteTableName, t.RemoteDbName, t.RemoteTableName)
		// 运行java命令
		cmd := exec.Command("java",
			"-cp", "mysql2mysql.jar;flink_libs/*",
			"com.demo.flink.FlinkCdcMySql",
			strCreatDatabase,
			strSource,
			strTarget,
			strInsert,
		)
		if err = cmd.Start(); err != nil {
			response.Fail(c, nil, "failed to start")
			return
		}
		strEs := fmt.Sprintf("CREATE TABLE flink_%s.flink_target_%s "+
			"%s"+ //主键必须标明NOT ENFORCED
			" WITH (\n"+
			"  'connector'  = 'elasticsearch-7',\n"+
			"  'hosts'     = 'http://47.120.73.205:9200',\n"+
			"  'index'        = '%d_%s_%s_index'\n"+
			")",
			projectTable2.PT_remote_db_name, projectTable2.PT_remote_table_name, configStr1,
			projectTable2.PU_uid, projectTable2.PT_remote_db_name, projectTable2.PT_remote_table_name)
		// 运行java命令
		cmd1 := exec.Command("java",
			"-cp", "mysql2es.jar;flink_libs/*",
			"com.demo.flink.FlinkCdcMySql",
			strCreatDatabase,
			strSource,
			strEs,
			strInsert,
		)
		if err = cmd1.Start(); err != nil {
			response.Fail(c, nil, "failed to start")
			return
		}
		// 保存Process以便后续控制
		process := cmd.Process
		process1 := cmd1.Process
		processMap.Store(projectTable2.PT_uid, [2]*os.Process{process, process1})

		response.Success(c, gin.H{"id": projectTable2.PT_uid}, "")
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}

func DeleteProjectTable(c *gin.Context) {
	type msg struct {
		Id uint `json:"id"`
	}

	var m msg
	if e := c.ShouldBindJSON(&m); e == nil {
		// 获取表
		var projectTable1 models.ProjectTable
		if e := mysql.DB.Where("PT_uid = ?", m.Id).First(&projectTable1).Error; e != nil {
			response.Fail(c, nil, "该项目表不存在!")
			return
		}
		if result := mysql.DB.Delete(&models.ProjectTable{}, m.Id); result.Error != nil {
			response.Fail(c, nil, "删除ProjectTable时出错")
			return
		} else {
			if result.RowsAffected == 0 {
				response.Success(c, nil, "要删除的记录不存在")
				return
			}
		}
		type ProcessList []*os.Process
		if processesInterface, ok := processMap.Load(m.Id); ok {
			// 确保加载的值是ProcessList类型
			if processes, ok := processesInterface.(ProcessList); ok {
				// 遍历数组，终止每个进程
				for _, proc := range processes {
					if err := proc.Kill(); err != nil {
						response.Fail(c, nil, "Failed to kill process1")
						return
					}
				}
				processMap.Delete(m.Id)
			}
			// 删除表
			dropTableSQL := fmt.Sprintf("DROP TABLE IF EXISTS %d_%s_%s;",
				projectTable1.PU_uid, projectTable1.PT_remote_db_name, projectTable1.PT_remote_table_name)
			if err := mysql.DB_flink.Exec(dropTableSQL).Error; err != nil {
				response.Fail(c, nil, "Failed to kill process")
				return
			}
			// 运行java命令
			cmd := exec.Command("java",
				"-cp", "deleteIndex.jar;flink_libs/*",
				"com.demo.flink.FlinkCdcMySql",
				fmt.Sprintf("/%d_%s_%s_index", projectTable1.PU_uid, projectTable1.PT_remote_db_name,
					projectTable1.PT_remote_table_name),
			)
			if err := cmd.Run(); err == nil {
				response.Success(c, nil, "Success to kill process and delete index")
				return
			}
		}
		response.Fail(c, nil, "Failed to kill process2")
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}

func GetProjectTable(c *gin.Context) {
	type table struct {
		Uid             uint   `json:"uid"`
		RemoteDbName    string `json:"remote_db_name"`
		RemoteTableName string `json:"remote_table_name"`
	}
	var t table

	if e := c.ShouldBindJSON(&t); e == nil {
		// 查询特定项目表
		var projectTable models.ProjectTable
		if e := mysql.DB.Where("PU_uid = ? and PT_remote_db_name = ? and PT_remote_table_name = ?",
			t.Uid, t.RemoteDbName, t.RemoteTableName).First(&projectTable).Error; e != nil {
			response.Fail(c, nil, "不存在该项目表!")
			return
		}
		// 获取表字段信息
		tmp, err := gorm.Open("mysql",
			fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				projectTable.PT_remote_username, projectTable.PT_remote_password,
				projectTable.PT_remote_hostname, projectTable.PT_remote_port, projectTable.PT_remote_db_name))
		if err != nil {
			response.Fail(c, nil, "数据库错误!")
			return
		}
		QueryColumnSQL := fmt.Sprintf("SELECT "+
			"COLUMN_NAME, DATA_TYPE, COLUMN_DEFAULT, IS_NULLABLE, COLUMN_KEY, COLUMN_COMMENT "+
			"FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '%s' "+
			"AND TABLE_NAME = '%s' ORDER BY ORDINAL_POSITION;",
			projectTable.PT_remote_db_name,
			projectTable.PT_remote_table_name)
		var columns []struct {
			ColumnName    string `gorm:"column:COLUMN_NAME"`
			DataType      string `gorm:"column:DATA_TYPE"`
			ColumnDefault string `gorm:"column:COLUMN_DEFAULT"`
			IsNullable    string `gorm:"column:IS_NULLABLE"`
			ColumnKey     string `gorm:"column:COLUMN_KEY"`
			ColumnComment string `gorm:"column:COLUMN_COMMENT"`
		}
		if err := tmp.Raw(QueryColumnSQL).Scan(&columns).Error; err != nil {
			response.Fail(c, nil, "数据库错误!")
			return
		}
		if err := tmp.Close(); err != nil {
			response.Fail(c, nil, "数据库错误!")
			return
		}
		type returnColumn struct {
			ColumnName    string `json:"name"`
			DataType      string `json:"data_type"`
			ColumnDefault string `json:"default"`
			IsNullable    string `json:"is_nullable"`
			ColumnKey     string `json:"key"`
			ColumnComment string `json:"comment"`
		}
		var returnColumns []returnColumn
		for _, column := range columns {
			returnColumns = append(returnColumns, returnColumn{column.ColumnName, column.DataType,
				column.ColumnDefault, column.IsNullable,
				column.ColumnKey, column.ColumnComment,
			})
		}
		response.Success(c, gin.H{"name": projectTable.PT_name, "description": projectTable.PT_description,
			"remote_hostname": projectTable.PT_remote_hostname, "remote_username": projectTable.PT_remote_username,
			"remote_password": projectTable.PT_remote_password, "remote_port": projectTable.PT_remote_port,
			"columns": returnColumns}, "")
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}

func GetAllProjectTable(c *gin.Context) {
	type table struct {
		Uid uint `json:"uid"`
	}
	var t table

	if e := c.ShouldBindJSON(&t); e == nil {
		// 查询特定项目表
		var projectTables []models.ProjectTable
		if e := mysql.DB.Where("PU_uid = ?", t.Uid).Find(&projectTables).Error; e != nil {
			response.Fail(c, nil, "不存在项目表!")
			return
		}
		type returnTable struct {
			RemoteDbName    string `json:"remote_db_name"`
			RemoteTableName string `json:"remote_table_name"`
		}
		var returnTables []returnTable
		for _, table1 := range projectTables {
			returnTables = append(returnTables, returnTable{
				table1.PT_remote_db_name,
				table1.PT_remote_table_name})
		}
		response.Success(c, gin.H{"tables": returnTables}, "")
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
