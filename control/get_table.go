package control

import (
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetTable(c *gin.Context) {
	type msg struct {
		Id uint `form:"id"`
	}
	var m msg
	if e := c.ShouldBindQuery(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}
	var tables []struct {
		TableName string `json:"table_name" gorm:"column:PT_name"`
		Id        uint   `json:"id" gorm:"column:PT_uid"`
	}
	if err := mysql.DB.Raw(`SELECT ProjectTable.PT_name,PT_uid FROM ProjectUser,ProjectTable 
	WHERE ProjectTable.PU_uid=ProjectUser.PU_uid AND ProjectUser.U_uid = ?`, m.Id).Scan(&tables).Error; err != nil {
		response.Fail(c, nil, "查询项目表时出错!")
		return
	}
	response.Success(c, gin.H{"tables": tables}, "")
}
