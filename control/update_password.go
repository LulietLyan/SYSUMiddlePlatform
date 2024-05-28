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
		//首先判断两次输入的密码是否一致
		if m.NewPwd != m.ConfirmPwd {
			response.Fail(c, nil, "两次输入的新密码不一致")
			return
		}
		//然后从token获取uid
		var userId uint
		if data, ok := c.Get("userId"); !ok {
			response.Fail(c, nil, "没有从token解析出所需信息")
			return
		} else {
			userId = data.(uint)
		}
		//根据uid找到用户记录
		tx := mysql.DB.Begin()
		var userRecord models.User
		if e := tx.Where("U_uid = ?", userId).First(&userRecord).Error; e != nil {
			tx.Rollback()
			response.Fail(c, nil, "找不到用户")
			return
		}
		//验证旧密码正确
		if userRecord.U_password != m.OldPwd {
			tx.Rollback()
			response.Fail(c, nil, "密码错误")
			return
		}
		//用新密码更新密码字段
		userRecord.U_password = m.NewPwd
		if e := tx.Save(&userRecord).Error; e != nil {
			tx.Rollback()
			response.Fail(c, nil, "更新密码时出错")
			return
		}
		//提交事务
		if e := tx.Commit().Error; e != nil {
			response.Fail(c, nil, "提交事务时出错")
			return
		}
		response.Success(c, nil, "")
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
