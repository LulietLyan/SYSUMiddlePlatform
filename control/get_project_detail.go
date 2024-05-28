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
		type column struct {
			ColumnName   string `json:"columnName"`
			ColumnType   string `json:"columnType"`
			IsPrimaryKey bool   `json:"isPrimaryKey"`
			IsForeignKey bool   `json:"isForeignKey"`
			IsNotNull    bool   `json:"isNotNull"`
			Desc         string `json:"desc"`
		}
		type table struct {
			TableName string   `json:"tableName"`
			TableDesc string   `json:"tableDesc"`
			Columns   []column `json:"columns"`
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
		var ptRecords []models.ProjectTable
		if e := mysql.DB.Where("PU_uid = ?", result.PU_uid).Find(&ptRecords).Error; e != nil {
			response.Fail(c, nil, "查找项目数据表时出错")
			return
		}
		var ts []table
		for _, ptRecord := range ptRecords {
			//查找ProjectColumn表
			var pcRecords []models.ProjectColumn
			if e := mysql.DB.Where("PT_uid = ?", ptRecord.PT_uid).Find(&pcRecords).Error; e != nil {
				response.Fail(c, nil, "查找项目数据列时出错")
				return
			}
			var cs []column
			for _, pcRecord := range pcRecords {
				cs = append(cs, column{ColumnName: pcRecord.PC_name, ColumnType: "暂无", IsNotNull: true, IsPrimaryKey: false, IsForeignKey: false, Desc: pcRecord.PC_description})
			}
			ts = append(ts, table{TableName: ptRecord.PT_name, TableDesc: ptRecord.PT_description, Columns: cs})
		}
		pd.Tables = ts
		response.Success(c, gin.H{"projectDetail": pd}, "")
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
