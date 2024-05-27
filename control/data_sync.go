package control

import (
	"backend/models"
	"backend/mysql"
	"fmt"
	"log"
	"os/exec"
)

func getAllConfig() (configList []*models.ProjectTable, err error) {
	if e := mysql.DB.Find(&configList).Error; e != nil {
		return nil, e
	}
	return configList, nil
}

func getSourceConfig(config *models.ProjectTable) (str string) {
	str = fmt.Sprintf("CREATE TABLE flink_%s.flink_source_%s ("+
		"%s"+ //主键必须标明NOT ENFORCED
		"NOT ENFORCED) WITH (\n"+
		"  'connector'  = 'mysql-cdc',\n"+
		"  'hostname'   = '%s',\n"+
		"  'port'   = '%d',\n"+
		"  'database-name'   = '%s',\n"+
		"  'table-name' = '%s',\n"+
		"  'username'   = '%s',\n"+
		"  'password'   = '%s'\n"+
		")",
		config.PT_remote_db_name, config.PT_remote_table_name, config.PT_config,
		config.PT_remote_hostname, config.PT_remote_port, config.PT_remote_db_name,
		config.PT_remote_table_name, config.PT_remote_username, config.PT_remote_password)
	return str
}

func getTargetConfig(config *models.ProjectTable) (str string) {
	str = fmt.Sprintf("CREATE TABLE flink_%s.flink_target_%s ("+
		"%s"+ //主键必须标明NOT ENFORCED
		"NOT ENFORCED) WITH (\n"+
		"  'connector'  = 'jdbc',\n"+
		"  'driver'     = 'com.mysql.cj.jdbc.Driver',\n"+
		"  'url'        = 'jdbc:mysql://47.121.29.57:3307/flink_target',\n"+
		"  'table-name' = '%d_%s_%s',\n"+
		"  'username'   = 'root',\n"+
		"  'password'   = '654321'\n"+
		")",
		config.PT_remote_db_name, config.PT_remote_table_name, config.PT_config,
		config.PU_uid, config.PT_remote_db_name, config.PT_remote_table_name)
	return str
}

func InitDataSync() error {
	todoList, err := getAllConfig()
	if err != nil {
		return err
	}
	// 运行所有配置信息
	for index, _ := range todoList {
		// 数据同步——创建数据表结构
		err = mysql.DB_flink.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %d_%s_%s;",
			todoList[index].PU_uid, todoList[index].PT_remote_db_name, todoList[index].PT_remote_table_name)).Error
		if err != nil {
			log.Fatalf("failed to drop table if exists: %v", err)
			return err
		}
		createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %d_%s_%s (%s);",
			todoList[index].PU_uid, todoList[index].PT_remote_db_name, todoList[index].PT_remote_table_name,
			todoList[index].PT_config)
		if err := mysql.DB_flink.Exec(createTableSQL).Error; err != nil {
			log.Fatalf("failed to create table: %v", err)
			return err
		}
		// 数据同步——java程序执行
		strCreatDatabase := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS flink_%s",
			todoList[index].PT_remote_db_name)
		strSource := getSourceConfig(todoList[index])
		strTarget := getTargetConfig(todoList[index])
		strInsert := fmt.Sprintf("INSERT INTO flink_%s.flink_target_%s select * from flink_%s.flink_source_%s\n",
			todoList[index].PT_remote_db_name, todoList[index].PT_remote_table_name,
			todoList[index].PT_remote_db_name, todoList[index].PT_remote_table_name)
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
