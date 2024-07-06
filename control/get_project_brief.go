package control

import (
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetProjectBrief(c *gin.Context) {
	var identity string
	var userId uint
	if data, ok := c.Get("identity"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
		return
	} else {
		identity = data.(string)
	}
	if data, ok := c.Get("userId"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
		return
	} else {
		userId = data.(uint)
	}
	var level int
	switch identity {
	case "Admin":
		level = 1
	case "Developer":
		level = 2
	default:
		level = 0
	}
	type msg struct {
		Offset uint `form:"offset"`
		Limit  uint `form:"limit"`
	}
	var m msg
	if e := c.ShouldBindQuery(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}

	var results []struct {
		U_username     string `gorm:"column:U_username;type:VARCHAR(64)" json:"U_username"`
		PU_description string `gorm:"column:PU_description;type:VARCHAR(8192)" json:"PU_description"`
		PU_logo_url    string `gorm:"column:PU_logo_url;type:VARCHAR(1024)" json:"PU_logo_url"`
		PU_uid         uint   `gorm:"column:PU_uid" json:"PU_uid"`
	}
	if e := mysql.DB.Table("User").
		Select("User.U_username, ProjectUser.PU_logo_url,ProjectUser.PU_description,ProjectUser.PU_uid").
		Joins("inner join ProjectUser on User.U_uid = ProjectUser.U_uid AND (1=? OR User.U_uid!=?) ", level, userId).Offset(m.Offset).Limit(m.Limit).
		Find(&results).Error; e != nil {
		response.Fail(c, nil, "查找项目时出错")
		return
	} else {
		type project struct {
			Title   string `json:"title"`
			Image   string `json:"image"`
			Content string `json:"content"`
			Id      uint   `json:"id"`
		}
		var projects []project
		for _, result := range results {
			projects = append(projects, project{Title: result.U_username, Image: result.PU_logo_url, Content: result.PU_description, Id: result.PU_uid})
		}
		response.Success(c, gin.H{"projects": projects}, "")
	}
}
