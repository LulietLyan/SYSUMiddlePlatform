package control

import (
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetProjectPageNumSearch(c *gin.Context) {
	type msg struct {
		Limit  uint   `json:"limit"`
		Search string `json:"search"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
	}

	var count int64
	search := "%" + m.Search + "%"
	if e := mysql.DB.Table("User").
		Joins("left join ProjectUser on User.U_uid = ProjectUser.U_uid").Where("User.U_username like ?", search).Count(&count).Error; e != nil {
		response.Fail(c, nil, "查找通知数量时出错")
		return
	} else {
		pages := count / int64(m.Limit)
		if count%int64(m.Limit) != 0 {
			pages++
		}
		response.Success(c, gin.H{"pages": pages}, "")
	}
}
