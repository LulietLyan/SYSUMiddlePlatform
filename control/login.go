package control

import (
	"backend/logic"
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	type msg struct {
		UserName string `json:"username"`
		Password string `json:"password"`
		Identity string `json:"identity"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e == nil {
		switch m.Identity {
		case "Admin":
			var adminRecord models.Admin
			if e = mysql.DB.Where("Admin_username = ?", m.UserName).First(&adminRecord).Error; e == nil {
				if adminRecord.Admin_password == m.Password {
					token, _ := logic.GenToken(adminRecord.Admin_uid, m.Identity)
					response.Success(c, gin.H{"token": token}, "")
				} else {
					response.Fail(c, nil, "密码错误")
				}
			} else {
				response.Fail(c, nil, "用户名不存在")
			}
		case "Analyzer":
			var auRecord models.AnalyticalUser
			if e = mysql.DB.Where("AU_username = ?", m.UserName).First(&auRecord).Error; e == nil {
				if auRecord.AU_password == m.Password {
					token, _ := logic.GenToken(auRecord.AU_uid, m.Identity)
					response.Success(c, gin.H{"token": token}, "")
				} else {
					response.Fail(c, nil, "密码错误")
				}
			} else {
				response.Fail(c, nil, "用户名不存在")
			}
		case "Developer":
			var puRecord models.ProjectUser
			if e = mysql.DB.Where("PU_username = ?", m.UserName).First(&puRecord).Error; e == nil {
				if puRecord.PU_password == m.Password {
					token, _ := logic.GenToken(puRecord.PU_uid, m.Identity)
					response.Success(c, gin.H{"token": token}, "")
				} else {
					response.Fail(c, nil, "密码错误")
				}
			} else {
				response.Fail(c, nil, "用户名不存在")
			}
		default:
			response.Fail(c, nil, "Identity参数为未知值")
		}
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
