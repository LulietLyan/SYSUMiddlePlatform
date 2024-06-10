package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func UpdatePermission(c *gin.Context) {
	//从上下文获取用户信息
	var identity string
	if data, ok := c.Get("identity"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
		return
	} else {
		identity = data.(string)
	}
	if identity != "Admin" {
		response.Fail(c, nil, "只有管理员可以使用此api!")
		return
	}
	//解析请求参数
	type msg struct {
		ProjectId uint   `json:"projectname"`
		TableId   uint   `json:"tablename"`
		Level     string `json:"level"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}
	//查找PU_uid
	var puRecord models.ProjectUser
	if e := mysql.DB.Where("U_uid = ?", m.ProjectId).First(&puRecord).Error; e != nil {
		response.Fail(c, nil, "找不到项目!")
		return
	}

	var pRecord models.Permission
	e := mysql.DB.Where("PU_uid = ? AND PT_uid = ?", puRecord.PU_uid, m.TableId).First(&pRecord).Error
	switch m.Level {
	case "Null":
		if e != nil {
			response.Success(c, nil, "用户本就无权限")
			return
		}
		if result := mysql.DB.Delete(&models.Permission{}, pRecord.P_uid); result.Error != nil {
			response.Fail(c, nil, "删除权限时出错")
			return
		} else {
			if result.RowsAffected == 0 {
				response.Success(c, nil, "要删除的权限不存在")
			} else {
				response.Success(c, nil, "")
			}
		}
	case "Read":
		pRecord.PU_uid = puRecord.PU_uid
		pRecord.PT_uid = m.TableId
		pRecord.P_level = 1
		if e := mysql.DB.Save(&pRecord).Error; e != nil {
			response.Fail(c, nil, "保存权限时出错")
			return
		}
		response.Success(c, nil, "")
	case "ReadWrite":
		pRecord.PU_uid = puRecord.PU_uid
		pRecord.PT_uid = m.TableId
		pRecord.P_level = 2
		if e := mysql.DB.Save(&pRecord).Error; e != nil {
			response.Fail(c, nil, "保存权限时出错")
			return
		}
		response.Success(c, nil, "")
	default:
		response.Fail(c, nil, "未知的权限等级")
		return
	}
}
