package control

import (
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetPermission(c *gin.Context) {
	//解析请求参数
	type msg struct {
		TableId uint `json:"tablename"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}
	var records []struct {
		ProjectName string `json:"projectname" gorm:"column:U_username"`
		Level       uint   `json:"level" gorm:"column:P_level"`
	}
	if e := mysql.DB.Raw(`
	SELECT User.U_username,Permission.P_level 
	FROM User,ProjectUser,Permission 
	WHERE User.U_uid=ProjectUser.U_uid AND ProjectUser.PU_uid=Permission.PU_uid AND Permission.PT_uid = ?
	`, m.TableId).Scan(&records).Error; e != nil {
		response.Success(c, nil, "")
		return
	}
	type auth struct {
		ProjectName string `json:"projectname" gorm:"column:PU_username"`
		Level       string `json:"level" gorm:"column:P_level"`
	}
	var auths []auth
	for _, record := range records {
		var levels = [2]string{"只读", "读写"}
		auths = append(auths, auth{ProjectName: record.ProjectName, Level: levels[record.Level-1]})
	}
	response.Success(c, gin.H{"auths": auths}, "")
}
