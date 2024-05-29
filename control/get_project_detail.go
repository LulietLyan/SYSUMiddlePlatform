package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetProjectDetail(c *gin.Context) {
	type msg struct {
		Projectname string `form:"projectname"`
	}
	var m msg
	if e := c.ShouldBindQuery(&m); e == nil {
		type member struct {
			Name  string `json:"name"`
			Phone string `json:"phone"`
			Email string `json:"email"`
			Job   string `json:"job"`
		}
		type table struct {
			RemoteDbName    string `json:"remote_db_name"`
			RemoteTableName string `json:"remote_table_name"`
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
		for _, pmRecords := range pmRecords {
			ms = append(ms, member{Name: pmRecords.PM_name, Phone: pmRecords.PM_phone, Email: pmRecords.PM_email, Job: pmRecords.PM_position})
		}
		pd.Members = ms
		//查找ProjectTable表
		var projectTables []models.ProjectTable
		if e := mysql.DB.Where("PU_uid = ?", result.PU_uid).Find(&projectTables).Error; e != nil {
			response.Fail(c, nil, "不存在项目表!")
			return
		}
		var returnTables []table
		for _, table1 := range projectTables {
			returnTables = append(returnTables, table{
				table1.PT_remote_db_name,
				table1.PT_remote_table_name})
		}
		pd.Tables = returnTables
		response.Success(c, gin.H{"projectDetail": pd}, "")
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
