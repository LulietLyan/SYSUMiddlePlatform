package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func DeleteApi(c *gin.Context) {
	type msg struct {
		Id uint `json:"id"`
	}

	var m msg
	if e := c.ShouldBindJSON(&m); e == nil {
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
		if identity != "Admin" { //不是管理员，只能删除自己的api
			var puRecord models.ProjectUser
			if e := mysql.DB.Where("U_uid = ?", userId).First(&puRecord).Error; e != nil {
				response.Fail(c, nil, "查找项目用户时失败")
				return
			}
			var apiRecord models.Api
			if e := mysql.DB.Where("A_uid = ?", m.Id).First(&apiRecord).Error; e != nil {
				response.Fail(c, nil, "查找api时失败")
				return
			}
			if puRecord.PU_uid != apiRecord.PU_uid {
				response.Fail(c, nil, "不能删除其他用户创建的api")
				return
			}
		}
		if result := mysql.DB.Delete(&models.Api{}, m.Id); result.Error != nil {
			response.Fail(c, nil, "删除Api时出错")
			return
		} else {
			if result.RowsAffected == 0 {
				response.Success(c, nil, "要删除的记录不存在")
			} else {
				response.Success(c, nil, "")
			}
		}
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
