package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

// 验证发来的用户名和密码（可能报告账号不存在或者密码错误），并生成一个token，随完整用户数据一起返回
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

		if acType[acRecord.AC_type-1] != m.Identity {
			var acTypeName = [3]string{"开发", "数据分析", "管理"}
			tx.Rollback()
			response.Fail(c, nil, "该邀请码只能用于创建"+acTypeName[acRecord.AC_type-1]+"用户")
			return
		}

		if m.Password1 != m.Password2 {
			tx.Rollback()
			response.Fail(c, nil, "两次输入的密码不一致")
		}

		switch m.Identity {
		case "Admin":
			//检查用户名不重复
			var adminRecord models.Admin
			if e := tx.Where("Admin_username=?", m.UserName).First(&adminRecord).Error; e == nil {
				tx.Rollback()
				response.Fail(c, nil, "用户名已存在")
				return
			}
			adminRecord = models.Admin{
				Admin_password: m.Password1,
				Admin_username: m.UserName,
			}
			if e := tx.Create(&adminRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "插入新用户信息时出错")
				return
			}
			acRecord.AC_usable = 0
			if e := tx.Save(&acRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "更新邀请码状态时出错")
				return
			}
			if e := tx.Commit().Error; e != nil {
				response.Fail(c, nil, "提交事务时出错")
				return
			}
			response.Success(c, gin.H{"username": m.UserName}, "")
			break
		case "Analyzer":
			var auRecord models.AnalyticalUser
			if e := tx.Where("AU_username=?", m.UserName).First(&auRecord).Error; e == nil {
				tx.Rollback()
				response.Fail(c, nil, "用户名已存在")
				return
			}
			auRecord = models.AnalyticalUser{
				AU_password: m.Password1,
				AU_username: m.UserName,
				AU_phone:    m.Phone,
				AU_std_uid:  m.Id,
				AU_email:    m.Email,
				AU_realname: m.Name,
			}
			if e := tx.Create(&auRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "插入新用户信息时出错")
				return
			}
			acRecord.AC_usable = 0
			if e := tx.Save(&acRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "更新邀请码状态时出错")
				return
			}
			if e := tx.Commit().Error; e != nil {
				response.Fail(c, nil, "提交事务时出错")
				return
			}
			response.Success(c, gin.H{"username": m.UserName}, "")
			break
		case "Developer":
			var puRecord models.ProjectUser
			if e := tx.Where("PU_username=?", m.UserName).First(&puRecord).Error; e == nil {
				tx.Rollback()
				response.Fail(c, nil, "用户名已存在")
				return
			}
			puRecord = models.ProjectUser{
				PU_password: m.Password1,
				PU_username: m.UserName,
				PU_email:    m.Email,
			}
			if e := tx.Create(&puRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "插入新用户信息时出错")
				return
			}
			acRecord.AC_usable = 0
			if e := tx.Save(&acRecord).Error; e != nil {
				tx.Rollback()
				response.Fail(c, nil, "更新邀请码状态时出错")
				return
			}
			if e := tx.Commit().Error; e != nil {
				response.Fail(c, nil, "提交事务时出错")
				return
			}
			response.Success(c, gin.H{"username": m.UserName}, "")
			break
		default:
			response.Fail(c, nil, "Identity参数为未知值")
		}
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
