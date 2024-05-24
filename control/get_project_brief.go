package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetProjectBrief(c *gin.Context) {
	var puRecords []models.ProjectUser
	if e := mysql.DB.Find(&puRecords).Error; e != nil {
		response.Fail(c, nil, "查找项目时出错")
		return
	} else {
		type project struct {
			Title   string `json:"title"`
			Image   string `json:"image"`
			Content string `json:"content"`
		}
		var projects []project
		for _, puRecord := range puRecords {
			projects = append(projects, project{Title: puRecord.PU_username, Image: puRecord.PU_logo_url, Content: puRecord.PU_description})
		}
		response.Success(c, gin.H{"projects": projects}, "")
	}
}
