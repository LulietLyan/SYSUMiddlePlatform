// package control

// import (
// 	"backend/models"
// 	"backend/mysql"
// 	"backend/response"

// 	"github.com/gin-gonic/gin"
// )

//	func GetProjectDetail(c *gin.Context) {
//		type msg struct {
//			Projectname string `form:"projectname"`
//		}
//		var m msg
//		if e := c.ShouldBindQuery(&m); e == nil {
//			type member struct {
//				Name  string `json:"name"`
//				Phone string `json:"phone"`
//				Email string `json:"email"`
//				Job   string `json:"job"`
//			}
//			type table struct {
//				RemoteDbName    string `json:"remote_db_name"`
//				RemoteTableName string `json:"remote_table_name"`
//			}
//			type projectDetail struct {
//				Logo        string   `json:"logo"`
//				ProjectName string   `json:"projectname"`
//				Description string   `json:"description"`
//				Members     []member `json:"members"`
//				Tables      []table  `json:"tables"`
//			}
//			//查找ProjectUser表
//			var result struct {
//				PU_logo_url    string `gorm:"column:PU_logo_url;type:VARCHAR(1024)" json:"PU_logo_url"`
//				U_username     string `gorm:"column:U_username;type:VARCHAR(64)" json:"U_username"`
//				PU_description string `gorm:"column:PU_description;type:VARCHAR(8192)" json:"PU_description"`
//				PU_uid         uint   `gorm:"column:PU_uid" json:"PU_uid"`
//			}
//			if e := mysql.DB.Table("User").Select("ProjectUser.PU_logo_url,User.U_username,ProjectUser.PU_description,ProjectUser.PU_uid").Joins("Join ProjectUser on User.U_uid = ProjectUser.U_uid").First(&result).Error; e != nil {
//				response.Fail(c, nil, "查找项目时出错")
//				return
//			}
//			var pd projectDetail
//			pd.Logo, pd.ProjectName, pd.Description = result.PU_logo_url, result.U_username, result.PU_description
//			//查找ProjectMember表
//			var pmRecords []models.ProjectMember
//			if e := mysql.DB.Where("PU_uid = ?", result.PU_uid).Find(&pmRecords).Error; e != nil {
//				response.Fail(c, nil, "查找项目成员时出错")
//				return
//			}
//			var ms []member
//			for _, pmRecords := range pmRecords {
//				ms = append(ms, member{Name: pmRecords.PM_name, Phone: pmRecords.PM_phone, Email: pmRecords.PM_email, Job: pmRecords.PM_position})
//			}
//			pd.Members = ms
//			//查找ProjectTable表
//			var projectTables []models.ProjectTable
//			if e := mysql.DB.Where("PU_uid = ?", result.PU_uid).Find(&projectTables).Error; e != nil {
//				response.Fail(c, nil, "不存在项目表!")
//				return
//			}
//			var returnTables []table
//			for _, table1 := range projectTables {
//				returnTables = append(returnTables, table{
//					table1.PT_remote_db_name,
//					table1.PT_remote_table_name})
//			}
//			pd.Tables = returnTables
//			response.Success(c, gin.H{"projectDetail": pd}, "")
//		} else { //JSON解析失败
//			response.Fail(c, nil, "数据格式错误!")
//		}
//	}
package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type column struct {
	ColumnName    string `json:"name"`
	DataType      string `json:"data_type"`
	ColumnDefault string `json:"default"`
	IsNullable    bool   `json:"is_nullable"`
	ColumnKey     string `json:"key"`
	ColumnComment string `json:"comment"`
}

func GetProjectDetail(c *gin.Context) {
	type msg struct {
		Projectname string `form:"projectname"`
	}
	var m msg
	if e := c.ShouldBindQuery(&m); e == nil {
		type member struct {
			Id    uint   `json:"id"`
			Name  string `json:"name"`
			Phone string `json:"phone"`
			Email string `json:"email"`
			Job   string `json:"job"`
		}

		type table struct {
			Id        uint     `json:"id"`
			TableName string   `json:"tableName"`
			TableDesc string   `json:"tableDesc"`
			Columns   []column `json:"columns"`
			Message   string   `json:"msg"`
			// RemoteDbName    string `json:"remote_db_name"`
			// RemoteTableName string `json:"remote_table_name"`
		}
		type projectDetail struct {
			Logo        string   `json:"logo"`
			ProjectName string   `json:"projectname"`
			Description string   `json:"description"`
			Members     []member `json:"members"`
			Tables      []table  `json:"tables"`
		}
		//查找ProjectUser表
		var result struct {
			PU_logo_url    string `gorm:"column:PU_logo_url;type:VARCHAR(1024)" json:"PU_logo_url"`
			U_username     string `gorm:"column:U_username;type:VARCHAR(64)" json:"U_username"`
			PU_description string `gorm:"column:PU_description;type:VARCHAR(8192)" json:"PU_description"`
			PU_uid         uint   `gorm:"column:PU_uid" json:"PU_uid"`
		}
		if e := mysql.DB.Table("User").Select("ProjectUser.PU_logo_url,User.U_username,ProjectUser.PU_description,ProjectUser.PU_uid").Joins("Join ProjectUser on User.U_uid = ProjectUser.U_uid").First(&result).Error; e != nil {
			response.Fail(c, nil, "查找项目时出错")
			return
		}
		var pd projectDetail
		pd.Logo, pd.ProjectName, pd.Description = result.PU_logo_url, result.U_username, result.PU_description
		//查找ProjectMember表
		var pmRecords []models.ProjectMember
		if e := mysql.DB.Where("PU_uid = ?", result.PU_uid).Find(&pmRecords).Error; e != nil {
			response.Fail(c, nil, "查找项目成员时出错")
			return
		}
		var ms []member
		for _, pmRecord := range pmRecords {
			ms = append(ms, member{Id: pmRecord.PM_uid, Name: pmRecord.PM_name, Phone: pmRecord.PM_phone, Email: pmRecord.PM_email, Job: pmRecord.PM_position})
		}
		pd.Members = ms
		//查找ProjectTable表
		var ptRecords []models.ProjectTable
		if e := mysql.DB.Where("PU_uid = ?", result.PU_uid).Find(&ptRecords).Error; e != nil {
			response.Fail(c, nil, "查找项目数据表时出错")
			// var projectTables []models.ProjectTable
			// if e := mysql.DB.Where("PU_uid = ?", result.PU_uid).Find(&projectTables).Error; e != nil {
			// 	response.Fail(c, nil, "不存在项目表!")
			return
		}
		var ts []table
		for _, ptRecord := range ptRecords {
			// 获取表字段信息
			cs, msg := queryColumn(ptRecord)
			ts = append(ts, table{TableName: ptRecord.PT_name, TableDesc: ptRecord.PT_description, Columns: cs, Message: msg})
		}
		pd.Tables = ts
		// pd.Tables = returnTables
		response.Success(c, gin.H{"projectDetail": pd}, "")
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}

func queryColumn(ptRecord models.ProjectTable) ([]column, string) {
	tmp, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			ptRecord.PT_remote_username, ptRecord.PT_remote_password,
			ptRecord.PT_remote_hostname, ptRecord.PT_remote_port, ptRecord.PT_remote_db_name))
	if err == nil {
		QueryColumnSQL := fmt.Sprintf("SELECT "+
			"COLUMN_NAME, DATA_TYPE, COLUMN_DEFAULT, IS_NULLABLE, COLUMN_KEY, COLUMN_COMMENT "+
			"FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '%s' "+
			"AND TABLE_NAME = '%s' ORDER BY ORDINAL_POSITION;",
			ptRecord.PT_remote_db_name,
			ptRecord.PT_remote_table_name)
		var cms []struct {
			ColumnName    string `gorm:"column:COLUMN_NAME"`
			DataType      string `gorm:"column:DATA_TYPE"`
			ColumnDefault string `gorm:"column:COLUMN_DEFAULT"`
			IsNullable    string `gorm:"column:IS_NULLABLE"`
			ColumnKey     string `gorm:"column:COLUMN_KEY"`
			ColumnComment string `gorm:"column:COLUMN_COMMENT"`
		}
		if err := tmp.Raw(QueryColumnSQL).Scan(&cms).Error; err != nil {
			return nil, "查找表字段出错"
		}
		if err := tmp.Close(); err != nil {
			return nil, "关闭数据库连接时出错"
		}
		var cs []column
		for _, cm := range cms {
			cs = append(cs, column{cm.ColumnName, cm.DataType,
				cm.ColumnDefault, cm.IsNullable == "YES",
				cm.ColumnKey, cm.ColumnComment,
			})
		}
		return cs, ""
	} else {
		return nil, "打开数据库连接失败"
	}
}
