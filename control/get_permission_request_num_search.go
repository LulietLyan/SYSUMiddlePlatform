package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetPermissionRequestNumSearch(c *gin.Context) {
	type msg struct {
		Limit     uint `form:"limit"`
		ProjectId uint `form:"projectname"`
		TableId   uint `form:"tablename"`
	}
	var m msg
	if e := c.ShouldBindQuery(&m); e == nil {
		if m.Limit == 0 {
			m.Limit = 15
		}
		var count int64
		q := mysql.DB.Model(&models.PermissionRequest{}).Where("PR_status = 1")
		if m.ProjectId != 0 {
			//查找PU_uid
			var puRecord models.ProjectUser
			if e := mysql.DB.Where("U_uid = ?", m.ProjectId).First(&puRecord).Error; e != nil {
				response.Fail(c, nil, "找不到项目!")
				return
			}
			q = q.Where("PU_uid = ?", puRecord.PU_uid)
		}
		if m.TableId != 0 {
			q = q.Where("PT_uid = ?", m.TableId)
		}
		if e := q.Count(&count).Error; e != nil {
			response.Fail(c, nil, "查找Api数量时出错")
			return
		} else {
			pages := count / int64(m.Limit)
			if count%int64(m.Limit) != 0 {
				pages++
			}
			response.Success(c, gin.H{"pages": pages}, "")
		}
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
