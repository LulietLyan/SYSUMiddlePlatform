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
		if data, ok := c.Get("identity"); !ok {
			response.Fail(c, nil, "没有从token解析出所需信息")
			return
		} else {
			identity = data.(string)
		}
		if identity != "Admin" {
			response.Fail(c, nil, "权限不足")
			return
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
