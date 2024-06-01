package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func DeleteMember(c *gin.Context) {
	type msg struct {
		Id uint `json:"id"`
	}

	var m msg
	if e := c.ShouldBindJSON(&m); e == nil {
		if result := mysql.DB.Delete(&models.ProjectMember{}, m.Id); result.Error != nil {
			response.Fail(c, nil, "删除成员时出错")
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
