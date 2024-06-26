package control

import (
	"backend/SQLParser"
	"backend/mysql"
	"backend/response"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	// ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓ 从 token 解析用户 id
	if data, ok := c.Get("pu_uid"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
		return
	} else {
		pu_uid = data.(uint)
		// ↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑ 解析完毕

		// ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓ 暂定前端需要发送以下信息
		type msg struct {
			ProjectName string `json:"projectName"`
			TableName   string `json:"tableName"`
			SqlCommand  string `json:"sqlCommand"`
		}
		var m msg
		if e := c.ShouldBindJSON(&m); e != nil {
			response.Fail(c, gin.H{"data": "请检查数据格式"}, "提交事务时出错")
			return
		} else {
			// ↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑ 暂定前端需要发送以上信息

			// ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓ 检查权限
			var authType struct {
				PT_uid  uint `gorm:"column:P_uid;primary_key" json:"P_uid"`
				P_level uint `gorm:"column:P_level" json:"P_level"`
			}

			err := mysql.DB.Begin().Raw(`
				SELECT Permission.PT_uid, Permission.P_level
				FROM Permission, ProjectTable 
				WHERE Permission.pu_uid=? AND Permission.PT_uid = ProjectTable.PT_uid`, pu_uid).First(&authType).Error
			// 用户必须具有写权限，否则毫无意义
			if err != nil && authType.P_level < 2 {
				response.Fail(c, gin.H{"data": "无权限"}, "检查权限时出错")
				return
			} else {
				// ↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑ 权限检查完毕

				m.TableName = SQLParser.SQLTreeGenerator(m.SqlCommand).Table.TableRefs.Left.Source.Name.O

				// ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓ 查找数据源参数
				var result struct {
					PT_uid               uint   `gorm:"column:PT_uid" json:"PT_uid"`
					PT_remote_db_name    string `gorm:"column:PT_remote_db_name" json:"PT_remote_db_name"`
					PT_remote_table_name string `gorm:"column:PT_remote_table_name" json:"PT_remote_table_name"`
					PT_remote_hostname   string `gorm:"column:PT_remote_hostname" json:"PT_remote_hostname"`
					PT_remote_userName   string `gorm:"column:PT_remote_userName" json:"PT_remote_userName"`
					PT_remote_password   string `gorm:"column:PT_remote_password" json:"PT_remote_password"`
					PT_remote_port       uint   `gorm:"column:PT_remote_port" json:"PT_remote_port"`
				}
				err = mysql.DB.Begin().Raw(`
					SELECT ProjectTable.PT_uid, ProjectTable.PT_remote_db_name, ProjectTable.PT_remote_table_name, ProjectTable.PT_remote_hostname, ProjectTable.PT_remote_userName , ProjectTable.PT_remote_password, ProjectTable.PT_remote_port 
					FROM User, ProjectUser, ProjectTable
					WHERE ProjectTable.PT_name=? AND User.U_uid = ProjectUser.U_uid AND ProjectTable.PU_uid=? AND ProjectUser.PU_uid = ProjectTable.PU_uid`, m.TableName, pu_uid).First(&result).Error

				if err != nil {
					response.Fail(c, gin.H{"data": "无相关项目"}, "查找表时出错")
					return
				} else {
					// ↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑ 查找数据源参数

					// ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓ 替换为后台数据库以及表名
					// 首先连接用户数据源
					DB_Origin, e := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", result.PT_remote_userName, result.PT_remote_password, result.PT_remote_hostname, result.PT_remote_port, result.PT_remote_db_name))

					println(result.PT_remote_userName, result.PT_remote_password, result.PT_remote_hostname, result.PT_remote_port, result.PT_remote_db_name)

					if e != nil {
						response.Fail(c, gin.H{"data": "连接数据源出错"}, "连接数据源出错")
						return
					} else {
						m.SqlCommand = strings.ReplaceAll(m.SqlCommand, m.ProjectName, result.PT_remote_db_name)
						m.SqlCommand = strings.ReplaceAll(m.SqlCommand, m.TableName, result.PT_remote_table_name)

						response.Success(c, gin.H{"rowsAffected": DB_Origin.Exec(m.SqlCommand).RowsAffected}, "")
					}
				}
			}
		}
	}
}

// SuperviseReadingAuth 并不是可用的路由，仅用于管理用户的读权限
func SuperviseReadingAuth(mysqlUser string, mysqlPassword string, targetTable string, grantOrRevoke bool) bool {
	if grantOrRevoke {
		err := mysql.DB_Authorize.Exec(`
		GRANT SELECT ON flink_target.? TO '?'@'%' IDENTIFIED BY '?'
	`, targetTable, mysqlUser, mysqlPassword).Error

		if err != nil {
			return false
		}

		return true
	} else {
		err := mysql.DB_Authorize.Exec(`
		REVOKE SELECT ON flink_target.? FROM '?'@'%' IDENTIFIED BY '?' IDENTIFIED BY '?'
	`, targetTable, mysqlUser, mysqlPassword).Error

		if err != nil {
			return false
		}

		return true
	}
}
