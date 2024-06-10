package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func GenInviteCode(c *gin.Context) {
	type msg struct {
		Type string `json:"type"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}
	var acType uint
	switch m.Type {
	case "Developer":
		acType = 1
	case "Analyzer":
		acType = 2
	case "Admin":
		acType = 3
	default:
		response.Fail(c, nil, "用户类型未知!")
		return
	}
	var inviteCode string
	tx := mysql.DB.Begin()
	for {
		inviteCode = RandomString(6)
		var acRecord models.ActivationCode
		if e := tx.Where("AC_code = ?", inviteCode).First(&acRecord).Error; e != nil {
			//保证邀请码不重复
			break
		}
	}
	acRecord := models.ActivationCode{AC_code: inviteCode, AC_usable: 1, AC_type: acType}
	if e := tx.Create(&acRecord).Error; e != nil {
		tx.Rollback()
		response.Fail(c, nil, "插入新邀请码信息时出错")
		return
	}
	if e := tx.Commit().Error; e != nil {
		response.Fail(c, nil, "提交事务时出错")
		return
	}
	response.Success(c, gin.H{"invitecode": inviteCode}, "")
}

// 定义包含大小写字母和数字的字符集
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 初始化随机种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

// 生成指定长度的随机字符串
func RandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
