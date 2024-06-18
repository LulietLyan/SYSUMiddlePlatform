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

	//查找原本的权限，以及原本的权限请求
	var pRecord models.Permission
	var prRecord models.PermissionRequest
	e1 := tx.Where("PU_uid = ? AND PT_uid = ?", pu_uid, result.PT_uid).First(&pRecord).Error
	e2 := tx.Where("PU_uid = ? AND PT_uid = ? AND PR_status = 1", pu_uid, result.PT_uid).First(&prRecord).Error
	if (e1 == nil && prLevel > pRecord.P_level) || e1 != nil { //如果原本有权限且申请提高权限，或者无权限时请求权限
		if e2 == nil && prRecord.PR_level != prLevel { //如果原本有请求，更改请求的权限
			prRecord.PR_level = prLevel
			if e := tx.Save(&prRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "更新密码时出错")
				return
			}
		} else { //如果原本无请求，插入请求
			prRecord.PU_uid, prRecord.PT_uid, prRecord.PR_level, prRecord.PR_status = pu_uid, result.PT_uid, prLevel, 1
			if e := tx.Create(&prRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "插入请求信息时出错")
				return
			}
		}
	} else if prLevel <= pRecord.P_level { //如果原本有权限，请求降权或同等权限(请求同等权限就是为了删除权限请求)
		//降低权限
		if prLevel != pRecord.P_level {
			if prLevel > 0 { //降低为只读
				pRecord.P_level = prLevel
				if e := tx.Save(&pRecord).Error; e != nil {
					tx.Rollback()
					response.Fail(c, nil, "降低权限失败")
					return
				}
			} else { //降低为无权限
				err = tx.Where("PU_uid = ? AND PT_uid = ?", pu_uid, result.PT_uid).Delete(&models.Permission{}).Error
				if err != nil {
					tx.Rollback()
					response.Fail(c, nil, "删除权限时出错")
					return
				}
			}
		}
		//删除权限请求
		if e2 == nil {
			err = tx.Where("PU_uid = ? AND PT_uid = ? AND PR_status = 1", pu_uid, result.PT_uid).Delete(&models.PermissionRequest{}).Error
			if err != nil {
				tx.Rollback()
				response.Fail(c, nil, "删除权限时出错")
				return
			}
		}
	}

	if e := tx.Commit().Error; e != nil {
		response.Fail(c, nil, "提交事务时出错")
		return
	}
	response.Success(c, nil, "")

}
