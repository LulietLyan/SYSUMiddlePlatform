package control

import (
	"backend/logic"
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	type msg struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e == nil {
		var userRecord models.User
		if e = mysql.DB.Where("U_username = ?", m.UserName).First(&userRecord).Error; e == nil {
			if userRecord.U_password == m.Password {
				var Identitys = [3]string{"Developer", "Analyzer", "Admin"}
				Identity := Identitys[userRecord.U_type-1]
				var pu_uid uint
				if Identity == "Developer" {
					var puRecord models.ProjectUser
					if e := mysql.DB.Where("U_uid = ?", userRecord.U_uid).First(&puRecord).Error; e != nil {
						response.Fail(c, nil, "查找项目用户时失败")
						return
					}
					pu_uid = puRecord.PU_uid
				}
				token, _ := logic.GenToken(userRecord.U_uid, Identity, pu_uid)
				response.Success(c, gin.H{"token": token, "identity": Identity}, "")
			} else {
				response.Fail(c, nil, "密码错误")
			}
		} else {
			response.Fail(c, nil, "用户名不存在")
		}
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
