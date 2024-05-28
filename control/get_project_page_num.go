package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetProjectPageNum(c *gin.Context) {
	type msg struct {
		Limit uint `form:"limit"`
	}
	var m msg
	if e := c.ShouldBindQuery(&m); e == nil {
		if m.Limit == 0 {
			m.Limit = 15
		}
		var count int64
		if e := mysql.DB.Model(&models.ProjectUser{}).Count(&count).Error; e != nil {
			response.Fail(c, nil, "查找通知数量时出错")
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
