package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	type msg struct {
		Id          string `json:"id"`
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		ProjectName string `json:"projectname"`
		Name        string `json:"name"`
		UserName    string `json:"username"`
		Password1   string `json:"password1"`
		Password2   string `json:"password2"`
		Invitecode  string `json:"invitecode"`
		Identity    string `json:"identity"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e == nil {
		//检查邀请码
		var acRecord models.ActivationCode
		tx := mysql.DB.Begin()
		e := tx.Where("AC_code=?", m.Invitecode).First(&acRecord).Error
		if e != nil {
			tx.Rollback()
			response.Fail(c, nil, "邀请码不存在")
			return
		}
		var acType = [3]string{"Developer", "Analyzer", "Admin"}

		if acRecord.AC_usable == 0 {
			tx.Rollback()
			response.Fail(c, nil, "邀请码已失效")
			return
		}
		uType := acRecord.AC_type - 1
		if acType[uType] != m.Identity {
			var acTypeName = [3]string{"开发", "数据分析", "管理"}
			tx.Rollback()
			response.Fail(c, nil, "该邀请码只能用于创建"+acTypeName[acRecord.AC_type-1]+"用户")
			return
		}
		//检查输入的密码
		if m.Password1 != m.Password2 {
			tx.Rollback()
			response.Fail(c, nil, "两次输入的密码不一致")
		}
		//检查没有重名的用户
		var userRecord models.User
		if e := tx.Where("U_username=?", m.UserName).First(&userRecord).Error; e == nil {
			tx.Rollback()
			response.Fail(c, nil, "用户名已存在")
			return
		}
		//向user表添加记录
		tmpUserName := generateUsername(5)
		tmpUserPwd := generatePassword(10)
		newMysqlUser(tmpUserName, tmpUserPwd)
		userRecord = models.User{U_username: m.UserName, U_password: m.Password1, U_type: uType + 1, U_mysqlUserName: tmpUserName, U_mysqlUserPwd: tmpUserPwd}
		if e := tx.Create(&userRecord).Error; e != nil {
			tx.Rollback()
			response.Fail(c, nil, "插入新用户信息时出错")
			return
		}
		//如果注册的用户是开发用户或分析用户，向相应表补充额外数据
		if uType == 1 {
			auRecord := models.AnalyticalUser{
				U_uid:       userRecord.U_uid,
				AU_phone:    m.Phone,
				AU_std_uid:  m.Id,
				AU_email:    m.Email,
				AU_realname: m.Name,
			}
			if e := tx.Create(&auRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "插入分析用户信息时出错")
				return
			}
		} else if uType == 0 {
			puRecord := models.ProjectUser{
				U_uid:       userRecord.U_uid,
				PU_email:    m.Email,
				PU_logo_url: "/logo/default",
			}
			if e := tx.Create(&puRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "插入项目用户信息时出错")
				return
			}
		}
		//将邀请码设置为失效
		acRecord.AC_usable = 0
		if e := tx.Save(&acRecord).Error; e != nil {
			tx.Rollback()
			response.Fail(c, nil, "更新邀请码状态时出错")
			return
		}
		//提交事务
		if e := tx.Commit().Error; e != nil {
			response.Fail(c, nil, "提交事务时出错")
			return
		}
		response.Success(c, gin.H{"username": m.UserName}, "")
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
func generateUsername(length int) string {
	// 初始化种子
	rand.Seed(time.Now().UnixNano())

	// 字符集
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generatePassword(length int) string {
	// 初始化种子
	rand.Seed(time.Now().UnixNano())

	// 字符集
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func newMysqlUser(user string, password string) {
	createUserSQL := "CREATE USER '" + user + "'@'localhost' IDENTIFIED BY '" + password + "';"

	err := mysql.DB.Exec(createUserSQL).Error
	if err != nil {
		fmt.Println("mysql 用户创建失败！")
	}
}
