package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetPermissionRequestNum(c *gin.Context) {
	type msg struct {
		Limit uint `form:"limit"`
	}
	var m msg
	if e := c.ShouldBindQuery(&m); e == nil {
		if m.Limit == 0 {
			m.Limit = 15
		}
		var count int64
		if e := mysql.DB.Model(&models.PermissionRequest{}).Where("PR_status = 1").Count(&count).Error; e != nil {
			response.Fail(c, nil, "查找权限请求数量时出错")
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
