package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func ApplyForTableAuth(c *gin.Context) {
	//获取token中的申请用户uid
	var pu_uid uint
	if data, ok := c.Get("pu_uid"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
		return
	} else {
		pu_uid = data.(uint)
	}
	//解析请求参数
	type msg struct {
		Projectname string `json:"projectname"`
		TableName   string `json:"tableName"`
		AuthType    string `json:"authType"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "提交事务时出错")
		return
	}
	//翻译权限等级
	var prLevel uint
	switch m.AuthType {
	case "只读":
		prLevel = 1
	case "读写":
		prLevel = 2
	case "无权限":
		prLevel = 0
	default:
		response.Fail(c, nil, "未知的权限等级")
		return
	}
	//查找目标表的主键
	tx := mysql.DB.Begin()
	var result struct {
		PT_uid uint `gorm:"column:PT_uid" json:"PT_uid"`
	}
	err := tx.Raw(`
		SELECT ProjectTable.PT_uid FROM 
			(SELECT ProjectUser.PU_uid FROM User LEFT JOIN ProjectUser ON User.U_uid=ProjectUser.U_uid WHERE User.U_username = ? ) T1
			LEFT JOIN ProjectTable on T1.PU_uid = ProjectTable.PU_uid	
		WHERE ProjectTable.PT_name=?
	`, m.Projectname, m.TableName).First(&result).Error
	if err != nil {
		tx.Rollback()
		response.Fail(c, nil, "查找表时出错")
		return
	}
	if prLevel > 0 {
		pr := models.PermissionRequest{PU_uid: pu_uid, PT_uid: result.PT_uid, PR_level: prLevel, PR_status: 1}
		if e := tx.Create(&pr).Error; e != nil {
			tx.Rollback()
			response.Fail(c, nil, "插入请求信息时出错")
			return
		}
	} else {
		err = tx.Where("PU_uid = ? AND PT_uid = ?", pu_uid, result.PT_uid).Delete(&models.Permission{}).Error
		if err != nil {
			tx.Rollback()
			response.Fail(c, nil, "删除权限时出错")
			return
		}
	}
	if e := tx.Commit().Error; e != nil {
		response.Fail(c, nil, "提交事务时出错")
		return
	}
	response.Success(c, nil, "")

}
