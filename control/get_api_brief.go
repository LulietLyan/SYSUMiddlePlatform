package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetApiBrief(c *gin.Context) {
	type msg struct {
		Offset uint `form:"offset"`
		Limit  uint `form:"limit"`
	}

	var m msg
	if e := c.ShouldBindQuery(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}

	if m.Limit == 0 {
		m.Limit = 15
	}
	var apiRecords []models.Api
	if e := mysql.DB.Order("updated_at DESC").Offset(m.Offset).Limit(m.Limit).Find(&apiRecords).Error; e != nil {
		response.Fail(c, nil, "查找Api时出错")
		return
	}
	type apiInfo struct {
		Title   string `json:"title"`
		Id      uint   `json:"id"`
		Content string `json:"content"`
		Type    string `json:"type"`
		Time    string `json:"time"`
	}
	var apiInfos []apiInfo
	var apiType = [2]string{"Midtable", "User"}
	for _, apiRecord := range apiRecords {
		apiInfos = append(apiInfos, apiInfo{Title: apiRecord.A_name, Id: apiRecord.A_uid, Content: apiRecord.A_description, Type: apiType[apiRecord.A_type-1], Time: apiRecord.CreatedAt.Format("2006-01-02 15:04")})
	}
	response.Success(c, gin.H{"apiInfos": apiInfos}, "")
}
