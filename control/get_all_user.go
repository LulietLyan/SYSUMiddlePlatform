package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetAllUser(c *gin.Context) {
	//从上下文获取用户信息
	var identity string
	if data, ok := c.Get("identity"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
		return
	} else {
		identity = data.(string)
	}
	if identity != "Admin" {
		response.Fail(c, nil, "只有管理员可以调用该api")
		return
	}
	//解析请求参数
	var uRecords []models.User
	if e := mysql.DB.Model(&models.User{}).Find(&uRecords).Error; e != nil {
		response.Fail(c, nil, "查找用户时出错")
		return
	}
	type user struct {
		Username string `json:"username"`
		Identity string `json:"identity"`
		Id       uint   `json:"id"`
	}

	var users []user
	var Identitys = [3]string{"Developer", "Analyzer", "Admin"}
	for _, uRecord := range uRecords {
		users = append(users, user{Username: uRecord.U_username, Identity: Identitys[uRecord.U_type-1], Id: uRecord.U_uid})
	}
	response.Success(c, gin.H{"users": users}, "")
}
