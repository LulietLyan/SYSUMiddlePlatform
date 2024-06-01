package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetApiBrief(c *gin.Context) {
	//从上下文获取用户信息
	var identity string
	var userId uint
	if data, ok := c.Get("identity"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
		return
	} else {
		identity = data.(string)
	}
	if data, ok := c.Get("userId"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
		return
	} else {
		userId = data.(uint)
	}
	//解析请求参数
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
	var apiTypes = [4]string{"Midtable", "Require", "User", "Me"}
	for _, apiRecord := range apiRecords {
		typeindex := apiRecord.A_type - 1 //1 2 3 --> 0 1 2
		if identity == "Developer" && apiRecord.A_type == 3 && apiRecord.PU_uid == userId {
			typeindex = 3 //标记为自己提供的api
		}
		apiInfos = append(apiInfos, apiInfo{Title: apiRecord.A_name, Id: apiRecord.A_uid, Content: apiRecord.A_description, Type: apiTypes[typeindex], Time: apiRecord.CreatedAt.Format("2006-01-02 15:04")})
	}
	response.Success(c, gin.H{"apiInfos": apiInfos}, "")
}
