package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetApiSearch(c *gin.Context) {
	type msg struct {
		Offset uint   `json:"offset"`
		Limit  uint   `json:"limit"`
		Search string `json:"search"`
		Type   string `json:"type"`
	}

	var m msg
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}
	if m.Limit == 0 {
		m.Limit = 15
	}
	var aType uint
	switch m.Type {
	case "Midtable":
		aType = 1
	case "User":
		aType = 2
	case "":
		aType = 0
	default:
		response.Fail(c, nil, "api类型未知")
		return
	}
	search := "%" + m.Search + "%"
	var apiRecords []models.Api
	if aType > 0 {
		if e := mysql.DB.Order("updated_at DESC").Offset(m.Offset).Limit(m.Limit).Where("A_name like ? AND A_type = ?", search, aType).Find(&apiRecords).Error; e != nil {
			response.Fail(c, nil, "查找Api时出错")
			return
		}
	} else {
		if e := mysql.DB.Order("updated_at DESC").Offset(m.Offset).Limit(m.Limit).Where("A_name like ?", search).Find(&apiRecords).Error; e != nil {
			response.Fail(c, nil, "查找Api时出错")
			return
		}
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
