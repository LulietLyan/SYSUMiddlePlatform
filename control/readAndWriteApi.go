package control

import (
	"backend/SQLParser"
	"backend/models"
	"backend/mysql"
	"backend/response"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
	tokenOfUser_1 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJMb2dpblRpbWUiOiIyMDI0LTA2LTI1VDIxOjQwOjM1LjMyMzM3NDMrMDg6MDAiLCJVc2VySWQiOjYsIklkZW50aXR5IjoiRGV2ZWxvcGVyIiwiUFVfdWlkIjozfQ.iGCAMDQil6OkM8Z1dZr-6PBgyGDa800WbezQ7ZHF90U"
	writeURL      = "https://127.0.0.1:2020/api/rNw/request/write"
)

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
		var m struct {
			SqlCommand string `json:"sqlCommand"`
		}
		if e := c.ShouldBindJSON(&m); e != nil {
			response.Fail(c, gin.H{"data": "请检查数据格式"}, "提交事务时出错")
			return
		} else {
			// ↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑ 暂定前端需要发送以上信息

			// ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓ 检查权限
			var authType struct {
				PT_uid  uint `gorm:"column:PT_uid"`
				P_level uint `gorm:"column:P_level"`
			}
			TableName := SQLParser.SQLTreeGenerator(m.SqlCommand).Table.TableRefs.Left.Source.Name.O
			err := mysql.DB.Begin().Raw(`
				SELECT Permission.PT_uid, Permission.P_level
				FROM Permission, ProjectTable 
				WHERE Permission.pu_uid=? AND Permission.PT_uid = ProjectTable.PT_uid AND ProjectTable.PT_name=?`, pu_uid, TableName).First(&authType).Error
			// 用户必须具有写权限，否则毫无意义
			fmt.Println(authType.P_level)
			if err != nil || authType.P_level < 2 {
				response.Fail(c, gin.H{"data": "无权限"}, "权限不足")
				return
			} else {
				// ↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑ 权限检查完毕

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
					FROM ProjectTable
					WHERE ProjectTable.PT_uid=?`, authType.PT_uid).First(&result).Error

				if err != nil {
					fmt.Println(authType.PT_uid)
					fmt.Println(err.Error())
					response.Fail(c, nil, "查找表时出错")
					return
				} else {
					// ↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑ 查找数据源参数

					// ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓ 替换为后台数据库以及表名
					// 首先连接用户数据源
					DbOrigin, e := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", result.PT_remote_userName, result.PT_remote_password, result.PT_remote_hostname, result.PT_remote_port, result.PT_remote_db_name))

					// println(result.PT_remote_userName, result.PT_remote_password, result.PT_remote_hostname, result.PT_remote_port, result.PT_remote_db_name)

					if e != nil {
						response.Fail(c, nil, "连接数据源出错")
						return
					} else {
						m.SqlCommand = strings.Replace(m.SqlCommand, TableName, result.PT_remote_table_name, 1)

						response.Success(c, gin.H{"rowsAffected": DbOrigin.Exec(m.SqlCommand).RowsAffected}, "")
					}
				}
			}
		}
	}
}

func InterpretUserReadingRequest(c *gin.Context) {

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

func TestWriting(c *gin.Context) {
	var m struct {
		IWantToTestWriting string `json:"iWantToTestWriting"`
	}

	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, gin.H{"data": "请检查数据格式"}, "提交事务时出错")
		return
	} else {
		var findStudent []models.Student
		var findStudentAgain []models.Student
		var maxStudentId int
		maxStudentId = 500
		mysql.DB_Demon.Where("student_id>500").Find(&findStudent)

		println("学生 id 大于 500 的目前有以下这些：")
		println("*******************************************************************************************************************************************************************************************************")
		for _, eachStudent := range findStudent {
			println("student_id: ", eachStudent.StudentId,
				"student_type: ", eachStudent.StudentType,
				"gender: ", eachStudent.Gender,
				"ethnicity: ", eachStudent.Ethnicity,
				"birth_str: ", eachStudent.BirthStr,
				"education_level: ", eachStudent.EducationLevel,
				"political_status: ", eachStudent.PoliticalStatus,
				"hometown: ", eachStudent.Hometown,
				"gaokao_score: ", eachStudent.GaokaoScore,
				"grade: ", eachStudent.Grade,
				"class: ", eachStudent.Class,
			)
			if int(eachStudent.StudentId) > maxStudentId {
				maxStudentId = int(eachStudent.StudentId)
			}
		}
		println("*******************************************************************************************************************************************************************************************************")

		var data = make(map[string]interface{})

		data["projectName"] = "1"
		data["tableName"] = "Student"
		data["sqlCommand"] = "INSERT INTO Student(student_id, student_type, gender, ethnicity, birth_str, education_level, political_status, hometown, gaokao_score, grade, class) VALUES " + "(" + strconv.Itoa(maxStudentId+1) + ", '境内生', '男', '汉族', '1999|01', '本科生', '群众', 'unknown', 1000, 2024, 1);"
		fmt.Println(data["sqlCommand"])
		client := &http.Client{}

		// 将请求参数编码为JSON格式
		jsonData, err := json.Marshal(data)
		if err != nil {
			response.Fail(c, gin.H{"data": "JSON编码失败"}, "JSON编码失败")
			return
		}

		request, err := http.NewRequest("POST", writeURL, bytes.NewBuffer(jsonData))
		if err != nil {
			response.Fail(c, gin.H{"data": "POST 请求失败"}, "POST 请求失败")
			return
		}
		request.Header.Set("Authorization", tokenOfUser_1)
		resp, err := client.Do(request)
		defer request.Body.Close()

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			response.Fail(c, gin.H{"data": "请求写数据时出错"}, "请求写数据时出错")
			return
		}
		println(content)

		mysql.DB_Demon.Where("student_id>500").Find(&findStudentAgain)
		println("插入一条数据后学生 id 大于 500 的目前有以下这些：")
		println("*******************************************************************************************************************************************************************************************************")
		for _, eachStudent := range findStudentAgain {
			println("student_id: ", eachStudent.StudentId,
				"student_type: ", eachStudent.StudentType,
				"gender: ", eachStudent.Gender,
				"ethnicity: ", eachStudent.Ethnicity,
				"birth_str: ", eachStudent.BirthStr,
				"education_level: ", eachStudent.EducationLevel,
				"political_status: ", eachStudent.PoliticalStatus,
				"hometown: ", eachStudent.Hometown,
				"gaokao_score: ", eachStudent.GaokaoScore,
				"grade: ", eachStudent.Grade,
				"class: ", eachStudent.Class,
			)
		}
		println("*******************************************************************************************************************************************************************************************************")
	}
}
