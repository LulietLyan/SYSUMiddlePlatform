package control

import (
	"backend/SQLParser"
	"backend/mysql"
	"backend/response"
	"database/sql"
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
	"unsafe"
)

func scanRows2map(rows *sql.Rows) []map[string]string {
	res := make([]map[string]string, 0)               //  定义结果 map
	colTypes, _ := rows.ColumnTypes()                 // 列信息
	var rowParam = make([]interface{}, len(colTypes)) // 传入到 rows.Scan 的参数 数组
	var rowValue = make([]interface{}, len(colTypes)) // 接收数据一行列的数组

	for i, colType := range colTypes {
		rowValue[i] = reflect.New(colType.ScanType())           // 跟据数据库参数类型，创建默认值 和类型
		rowParam[i] = reflect.ValueOf(&rowValue[i]).Interface() // 跟据接收的数据的类型反射出值的地址

	}
	// 遍历
	for rows.Next() {
		rows.Scan(rowParam...) // 赋值到 rowValue 中
		record := make(map[string]string)
		for i, colType := range colTypes {

			if rowValue[i] == nil {
				record[colType.Name()] = ""
			} else {
				record[colType.Name()] = Byte2Str(rowValue[i].([]byte))
			}
		}
		res = append(res, record)
	}
	return res
}

// []byte to string
func Byte2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func InterpretUserWritingRequest(c *gin.Context) {
	var pu_uid uint

	// ****************************************** 从 token 解析用户 id
	if data, ok := c.Get("pu_uid"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
	} else {
		pu_uid = data.(uint)
	}
	// ****************************************** 解析完毕

	var m struct {
		projectName string `json:"projectName"`
		tableName   string `json:"tableName"`
		sqlCommand  string `json:"sqLCommand"`
	}
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "提交事务时出错")
	}

	// ****************************************** 检查权限
	var authType struct {
		PT_uid  uint `gorm:"column:P_uid;primary_key" json:"P_uid"`
		P_level uint `gorm:"column:P_level" json:"P_level"`
	}
	err := mysql.DB.Select(`
		SELECT Permission.P_level, Permission.PT_uid 
		From Permission, ProjectTable 
		Where Permission.pu_uid = ? AND Permission.PT_uid = ProjectTable.PT_uid`, pu_uid).First(&authType).Error
	if err != nil && authType.P_level < 2 {
		response.Fail(c, nil, "检查权限时出错")
	}
	// ****************************************** 权限检查完毕

	m.tableName = SQLParser.SQLTreeGenerator(m.sqlCommand).StmtTree.Table.Name.Original

	//查找目标表的主键
	tx := mysql.DB.Begin()
	var result struct {
		PT_uid               uint   `gorm:"column:PT_uid" json:"PT_uid"`
		PT_remote_db_name    string `gorm:"column:PT_remote_db_name;type:VARCHAR(64)" json:"PT_remote_db_name"`
		PT_remote_table_name string `gorm:"column:PT_remote_table_name;type:VARCHAR(64)" json:"PT_remote_table_name"`
	}
	err = tx.Raw(`
		SELECT ProjectTable.PT_uid, ProjectTable.PT_remote_db_name, ProjectTable.PT_remote_table_name FROM 
			(SELECT ProjectUser.PU_uid FROM User LEFT JOIN ProjectUser ON User.U_uid=ProjectUser.U_uid WHERE User.U_username = ? ) T1
			LEFT JOIN ProjectTable on T1.PU_uid = ProjectTable.PU_uid	
		WHERE ProjectTable.PT_name=?
	`, m.projectName, m.tableName).First(&result).Error
	if err != nil {
		tx.Rollback()
		response.Fail(c, nil, "查找表时出错")
	}

	// ****************************************** 替换为后台数据库以及表名
	m.sqlCommand = strings.ReplaceAll(m.sqlCommand, m.projectName, result.PT_remote_db_name)
	m.sqlCommand = strings.ReplaceAll(m.sqlCommand, m.tableName, result.PT_remote_table_name)

	rowsAffected := tx.Raw(m.sqlCommand).RowsAffected

	response.Success(c, gin.H{"rowsAffected": rowsAffected}, "")
}

func InterpretUserReadingRequest(c *gin.Context) {
	var pu_uid uint

	// ****************************************** 从 token 解析用户 id
	if data, ok := c.Get("pu_uid"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
	} else {
		pu_uid = data.(uint)
	}
	// ****************************************** 解析完毕

	var m struct {
		projectName string `json:"projectName"`
		tableName   string `json:"tableName"`
		sqlCommand  string `json:"sqLCommand"`
	}
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "提交事务时出错")
	}

	// ****************************************** 检查权限
	var authType struct {
		PT_uid  uint `gorm:"column:P_uid;primary_key" json:"P_uid"`
		P_level uint `gorm:"column:P_level" json:"P_level"`
	}
	err := mysql.DB.Select(`
		SELECT Permission.P_level, Permission.PT_uid 
		From Permission, ProjectTable 
		Where Permission.pu_uid = ? AND Permission.PT_uid = ProjectTable.PT_uid`, pu_uid).First(&authType).Error
	if err != nil || authType.PT_uid == 0 {
		response.Fail(c, nil, "检查权限时出错")
	}
	// ****************************************** 权限检查完毕

	response.Success(c, gin.H{}, "")
}
