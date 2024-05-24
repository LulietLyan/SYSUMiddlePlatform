package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func UpdatePassword(c *gin.Context) {
	type msg struct {
		OldPwd     string `json:"oldpwd"`
		NewPwd     string `json:"newpwd"`
		ConfirmPwd string `json:"confirmPwd"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e == nil {

		if m.NewPwd != m.ConfirmPwd {
			response.Fail(c, nil, "两次输入的新密码不一致")
			return
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

		tx := mysql.DB.Begin()
		switch identity {
		case "Admin":
			var adminRecord models.Admin
			if e := tx.Where("Admin_uid = ?", userId).First(&adminRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "找不到用户")
				return
			}
			if adminRecord.Admin_password != m.OldPwd {
				tx.Rollback()
				response.Fail(c, nil, "密码错误")
				return
			}
			adminRecord.Admin_password = m.NewPwd
			if e := tx.Save(&adminRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "更新密码时出错")
				return
			}
			if e := tx.Commit().Error; e != nil {
				response.Fail(c, nil, "提交事务时出错")
				return
			}
			response.Success(c, nil, "")
		case "Analyzer":
			var auRecord models.AnalyticalUser
			if e := tx.Where("AU_uid = ?", userId).First(&auRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "找不到用户")
				return
			}
			if auRecord.AU_password != m.OldPwd {
				tx.Rollback()
				response.Fail(c, nil, "密码错误")
				return
			}
			auRecord.AU_password = m.NewPwd
			if e := tx.Save(&auRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "更新密码时出错")
				return
			}
			if e := tx.Commit().Error; e != nil {
				response.Fail(c, nil, "提交事务时出错")
				return
			}
			response.Success(c, nil, "")
		case "Developer":
			var puRecord models.ProjectUser
			if e := tx.Where("PU_uid = ?", userId).First(&puRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "找不到用户")
				return
			}
			if puRecord.PU_password != m.OldPwd {
				tx.Rollback()
				response.Fail(c, nil, "密码错误")
				return
			}
			puRecord.PU_password = m.NewPwd
			if e := tx.Save(&puRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "更新密码时出错")
				return
			}
			if e := tx.Commit().Error; e != nil {
				response.Fail(c, nil, "提交事务时出错")
				return
			}
			response.Success(c, nil, "")
		default:
			response.Fail(c, nil, "Identity参数为未知值")
		}
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
