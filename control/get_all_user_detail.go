package control

import (
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetAllUserDetail(c *gin.Context) {
	//从上下文获取用户信息
	var identity string
	if data, ok := c.Get("identity"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
		return
	} else {
		identity = data.(string)
	}
	if identity != "Admin" {
		response.Fail(c, nil, "只有管理员用户可以调用此api!")
		return
	}
	var pUsers []struct {
		ProjectName  string `json:"projectname" gorm:"column:U_username"`
		Password     string `json:"password" gorm:"column:U_password"`
		NumOfMembers uint   `json:"num_of_members" gorm:"column:member_count"`
		NumOfTables  uint   `json:"num_of_tables" gorm:"column:table_count"`
	}
	if err := mysql.DB.Raw(`SELECT A.U_username,A.U_password,COALESCE(B.member_count,0) AS member_count,COALESCE(C.table_count,0) AS table_count
		FROM
		(SELECT User.U_uid,User.U_username,User.U_password,ProjectUser.PU_uid FROM User left join ProjectUser on User.U_uid = ProjectUser.U_uid
		WHERE User.U_type=1) A
		LEFT JOIN
		(SELECT PU_uid,COUNT(*) AS member_count FROM ProjectMember GROUP BY PU_uid) B ON A.PU_uid=B.PU_uid
		LEFT JOIN
		(SELECT PU_uid,COUNT(*) AS table_count FROM ProjectTable GROUP BY PU_uid) C ON A.PU_uid=C.PU_uid`).Scan(&pUsers).Error; err != nil {
		response.Fail(c, nil, "查询项目用户时出错!")
		return
	}
	var aUsers []struct {
		Username string `json:"username" gorm:"column:U_username"`
		Password string `json:"password" gorm:"column:U_password"`
		Name     string `json:"name" gorm:"column:AU_realname"`
		Uid      string `json:"uid" gorm:"column:AU_std_uid"`
		Phone    string `json:"phone" gorm:"column:AU_phone"`
		Email    string `json:"email" gorm:"column:AU_email"`
	}
	if err := mysql.DB.Raw(`SELECT U_username,U_password,AU_realname,AU_std_uid,AU_phone,AU_email
	FROM User LEFT JOIN AnalyticalUser ON User.U_uid = AnalyticalUser.U_uid
	WHERE User.U_type=2`).Scan(&aUsers).Error; err != nil {
		response.Fail(c, nil, "查询分析用户时出错!")
		return
	}
	response.Success(c, gin.H{"developers": pUsers, "analyzers": aUsers}, "")
}
