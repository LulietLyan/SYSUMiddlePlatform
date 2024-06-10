package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func ApprovePermissionRequest(c *gin.Context) {
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
		Id uint `json:"id"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}
	tx := mysql.DB.Begin()
	var prRecord models.PermissionRequest
	if e := tx.Where("PR_uid = ?", m.Id).First(&prRecord).Error; e != nil {
		tx.Rollback()
		response.Fail(c, nil, "找不到权限请求!")
		return
	}
	prRecord.PR_status = 2
	if e := tx.Save(&prRecord).Error; e != nil {
		println(e)
		tx.Rollback()
		response.Fail(c, nil, "更新请求状态时出错!")
		return
	}
	var pRecord models.Permission
	tx.Where("PU_uid = ? AND PT_uid = ?", prRecord.PU_uid, prRecord.PT_uid).First(&pRecord)
	switch prRecord.PR_level {
	case 1:
		pRecord.PU_uid = prRecord.PU_uid
		pRecord.PT_uid = prRecord.PT_uid
		pRecord.P_level = 1
		if e := tx.Save(&pRecord).Error; e != nil {
			tx.Rollback()
			response.Fail(c, nil, "保存权限时出错")
			return
		}
		if e := tx.Commit().Error; e != nil {
			response.Fail(c, nil, "提交事务时出错")
			return
		}
		response.Success(c, nil, "")
	case 2:
		pRecord.PU_uid = prRecord.PU_uid
		pRecord.PT_uid = prRecord.PT_uid
		pRecord.P_level = 2
		if e := tx.Save(&pRecord).Error; e != nil {
			tx.Rollback()
			response.Fail(c, nil, "保存权限时出错")
			return
		}
		if e := tx.Commit().Error; e != nil {
			response.Fail(c, nil, "提交事务时出错")
			return
		}
		response.Success(c, nil, "")
	default:
		tx.Rollback()
		response.Fail(c, nil, "请求的权限等级未知")
		return
	}
}
