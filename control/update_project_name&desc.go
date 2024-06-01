package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func UpdateProjectNameDesc(c *gin.Context) {
	//从上下文获取用户信息
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
	if identity != "Developer" {
		response.Fail(c, nil, "只有开发用户可以使用此api!")
		return
	}
	//解析请求参数
	type msg struct {
		Projectname string `json:"projectname"`
		Description string `json:"description"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}
	tx := mysql.DB.Begin()
	var uRecord models.User
	if e := tx.Where("U_username = ?", m.Projectname).First(&uRecord).Error; e == nil {
		if uRecord.U_uid != userId {
			tx.Rollback()
			response.Fail(c, nil, "已经有使用该名称的其他用户!")
			return
		}
	}
	uRecord = models.User{}
	if e := tx.Where("U_uid = ?", userId).First(&uRecord).Error; e != nil {
		tx.Rollback()
		response.Fail(c, nil, "找不到用户!")
		return
	}
	uRecord.U_username = m.Projectname
	if e := tx.Save(&uRecord).Error; e != nil {
		tx.Rollback()
		response.Fail(c, nil, "更新用户名时出错")
		return
	}
	var puRecord models.ProjectUser
	if e := mysql.DB.Raw("SELECT ProjectUser.* FROM User,ProjectUser WHERE User.U_uid=ProjectUser.U_uid and User.U_uid = ?", userId).First(&puRecord).Error; e != nil {
		response.Fail(c, nil, "查找项目记录时出错")
		return
	}
	puRecord.PU_description = m.Description
	if e := tx.Save(&puRecord).Error; e != nil {
		tx.Rollback()
		response.Fail(c, nil, "更新项目简介时出错")
		return
	}
	if e := tx.Commit().Error; e != nil {
		response.Fail(c, nil, "提交事务时出错")
		return
	}
	response.Success(c, nil, "")
}
