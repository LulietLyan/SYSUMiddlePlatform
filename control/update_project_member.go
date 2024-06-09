package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func UpdateProjectMember(c *gin.Context) {
	//从上下文获取用户信息
	var identity string
	var pu_uid uint
	if data, ok := c.Get("identity"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
		return
	} else {
		identity = data.(string)
	}
	if data, ok := c.Get("pu_uid"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
		return
	} else {
		pu_uid = data.(uint)
	}
	if identity != "Developer" {
		response.Fail(c, nil, "只有开发用户可以使用此api!")
		return
	}
	//解析请求参数
	type msg struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		Job   string `json:"job"`
		Id    int64  `json:"id"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}
	if m.Id < 0 {
		pmRecord := models.ProjectMember{PM_name: m.Name, PM_phone: m.Phone, PM_email: m.Email, PM_position: m.Job, PU_uid: pu_uid}
		if e := mysql.DB.Create(&pmRecord).Error; e != nil {
			response.Fail(c, nil, "插入新成员信息时出错")
			return
		}
		response.Success(c, nil, "")
		return
	} else {
		pmRecord := models.ProjectMember{PM_uid: uint(m.Id), PM_name: m.Name, PM_phone: m.Phone, PM_email: m.Email, PM_position: m.Job, PU_uid: pu_uid}
		if e := mysql.DB.Save(&pmRecord).Error; e != nil {
			response.Fail(c, nil, "更新成员信息时出错")
			return
		}
		response.Success(c, nil, "")
		return
	}
}
