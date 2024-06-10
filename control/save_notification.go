package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func SaveNotification(c *gin.Context) {
	//从上下文获取用户信息
	var identity string
	if data, ok := c.Get("identity"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
		return
	} else {
		identity = data.(string)
	}
	if identity != "Admin" {
		response.Fail(c, nil, "只有管理员可以调用该接口")
		return
	}
	//解析参数
	type msg struct {
		Type    string `json:"type"`
		Users   []uint `json:"users"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}
	switch m.Type {
	case "AllUser":
		nRecord := models.Notifications{N_type: 3, N_Title: m.Title, N_Body: m.Content}
		if e := mysql.DB.Create(&nRecord).Error; e != nil {
			response.Fail(c, nil, "插入新Api信息时出错")
			return
		}
		response.Success(c, nil, "")
	case "AllAnalyzer":
		nRecord := models.Notifications{N_type: 2, N_Title: m.Title, N_Body: m.Content}
		if e := mysql.DB.Create(&nRecord).Error; e != nil {
			response.Fail(c, nil, "插入新Api信息时出错")
			return
		}
		response.Success(c, nil, "")
	case "AllDeveloper":
		nRecord := models.Notifications{N_type: 1, N_Title: m.Title, N_Body: m.Content}
		if e := mysql.DB.Create(&nRecord).Error; e != nil {
			response.Fail(c, nil, "插入新Api信息时出错")
			return
		}
		response.Success(c, nil, "")
	case "Analyzer":
		tx := mysql.DB.Begin()
		for _, id := range m.Users {
			nRecord := models.Notifications{N_type: 5, AU_uid: id, N_Title: m.Title, N_Body: m.Content}
			if e := tx.Create(&nRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "插入新通知信息时出错")
				return
			}
			if e := tx.Commit().Error; e != nil {
				response.Fail(c, nil, "提交事务时出错")
				return
			}
			response.Success(c, nil, "")
		}
	case "Developer":
		tx := mysql.DB.Begin()
		for _, id := range m.Users {
			nRecord := models.Notifications{N_type: 4, PU_uid: id, N_Title: m.Title, N_Body: m.Content}
			if e := tx.Create(&nRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "插入新通知信息时出错")
				return
			}
			if e := tx.Commit().Error; e != nil {
				response.Fail(c, nil, "提交事务时出错")
				return
			}
			response.Success(c, nil, "")
		}
	default:
		response.Fail(c, nil, "未知的通知类型")
		return
	}
}
