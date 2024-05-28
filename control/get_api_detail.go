package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetApiDetail(c *gin.Context) {
	type msg struct {
		Id uint `json:"id"`
	}

	var m msg
	if e := c.ShouldBindQuery(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}

	var apiRecord models.Api
	if e := mysql.DB.Order("updated_at DESC").Where("A_uid = ?", m.Id).First(&apiRecord).Error; e != nil {
		response.Fail(c, nil, "查找Api时出错")
		return
	}
	type apiInfo struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Url      string `json:"url"`
		Desc     string `json:"desc"`
		Request  string `json:"request"`
		Response string `json:"response"`
		Id       uint   `json:"id"`
	}
	var apiType = [2]string{"Midtable", "User"}
	data := apiInfo{Name: apiRecord.A_name, Type: apiType[apiRecord.A_type-1], Url: apiRecord.A_url, Desc: apiRecord.A_description, Request: apiRecord.A_parameter, Response: apiRecord.A_respond, Id: apiRecord.A_uid}
	response.Success(c, gin.H{"apiInfo": data}, "")
}
