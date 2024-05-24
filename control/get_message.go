package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

// 验证发来的用户名和密码（可能报告账号不存在或者密码错误），并生成一个token，随完整用户数据一起返回
func GetMessage(c *gin.Context) {
	type msg struct {
		Offset uint `json:"offset"`
		Limit  uint `json:"limit"`
	}

	var m msg
	if e := c.ShouldBindJSON(&m); e == nil {
		if m.Limit == 0 {
			m.Limit = 15
		}

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

		var nRecords []models.Notifications
		switch identity {
		case "Admin":
			if e := mysql.DB.Order("UpdatedAt DESC").Offset(m.Offset).Limit(m.Limit).Find(&nRecords).Error; e != nil {
				response.Fail(c, nil, "查找通知时出错")
				return
			}
		case "Analyzer":
			if e := mysql.DB.Order("UpdateAt DESC").Where("N_type = 2 OR N_type = 3 OR (N_type = 5 AND AU_uid = ?)", userId).Offset(m.Offset).Limit(m.Limit).Find(&nRecords).Error; e != nil {
				response.Fail(c, nil, "查找通知时出错")
				return
			}
		case "Developer":
			if e := mysql.DB.Order("UpdateAt DESC").Where("N_type = 1 OR N_type = 3 OR (N_type = 4 AND PU_uid = ?)", userId).Offset(m.Offset).Limit(m.Limit).Find(&nRecords).Error; e != nil {
				response.Fail(c, nil, "查找通知时出错")
				return
			}
		default:
			response.Fail(c, nil, "Identity参数为未知值")
			return
		}

		type message struct {
			Id      uint   `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
			Author  string `json:"author"`
			Time    string `json:"time"`
		}
		var messages []message
		for _, mRecord := range nRecords {
			messages = append(messages, message{Id: mRecord.N_uid, Title: mRecord.N_Title, Content: mRecord.N_Body, Author: "数据中台管理团队", Time: mRecord.UpdatedAt.Format("2006-01-02 15:04")})
		}
		response.Success(c, gin.H{"messages": messages}, "")

	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
