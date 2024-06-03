package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetPermission(c *gin.Context) {
	//解析请求参数
	type msg struct {
		ProjectId uint `json:"projectname"`
		TableId   uint `json:"tablename"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}
	var pRecord models.Permission
	if e := mysql.DB.Where("PU_uid = ? AND PT_uid = ?", m.ProjectId, m.TableId).First(&pRecord).Error; e != nil {
		response.Success(c, nil, "")
		return
	}

}
