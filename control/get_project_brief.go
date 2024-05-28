package control

import (
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetProjectBrief(c *gin.Context) {
	var results []struct {
		U_username     string `gorm:"column:U_username;type:VARCHAR(64)" json:"U_username"`
		PU_description string `gorm:"column:PU_description;type:VARCHAR(8192)" json:"PU_description"`
		PU_logo_url    string `gorm:"column:PU_logo_url;type:VARCHAR(1024)" json:"PU_logo_url"`
	}
	if e := mysql.DB.Table("User").
		Select("User.U_username, ProjectUser.PU_logo_url,ProjectUser.PU_description").
		Joins("left join ProjectUser on User.U_uid = ProjectUser.U_uid").
		Find(&results).Error; e != nil {
		response.Fail(c, nil, "查找项目时出错")
		return
	} else {
		type project struct {
			Title   string `json:"title"`
			Image   string `json:"image"`
			Content string `json:"content"`
		}
		var projects []project
		for _, result := range results {
			projects = append(projects, project{Title: result.U_username, Image: result.PU_logo_url, Content: result.PU_description})
		}
		response.Success(c, gin.H{"projects": projects}, "")
	}
}
