package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetApiSearch(c *gin.Context) {
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
	search := "%" + m.Search + "%"
	//搜索api记录
	var apiRecords []models.Api
	switch m.Type {
	case "Midtable":
		if e := mysql.DB.Order("updated_at DESC").Offset(m.Offset).Limit(m.Limit).Where("A_name like ? AND A_type = ?", search, 1).Find(&apiRecords).Error; e != nil {
			response.Fail(c, nil, "查找Api时出错")
			return
		}
	case "Require":
		if e := mysql.DB.Order("updated_at DESC").Offset(m.Offset).Limit(m.Limit).Where("A_name like ? AND A_type = ?", search, 2).Find(&apiRecords).Error; e != nil {
			response.Fail(c, nil, "查找Api时出错")
			return
		}
	case "User":
		//其他用户提供的api
		if identity == "Admin" { //管理员查找时，User包括所有用户提供的api
			if e := mysql.DB.Order("updated_at DESC").Offset(m.Offset).Limit(m.Limit).Where("A_name like ? AND A_type = ?", search, 3).Find(&apiRecords).Error; e != nil {
				response.Fail(c, nil, "查找Api时出错")
				return
			}
		} else { //项目用户查找时，User只包括其他用户提供的api
			if e := mysql.DB.Order("updated_at DESC").Offset(m.Offset).Limit(m.Limit).Where("A_name like ? AND A_type = ? AND PU_uid != ?", search, 3, userId).Find(&apiRecords).Error; e != nil {
				response.Fail(c, nil, "查找Api时出错")
				return
			}
		}
	case "Me":
		//只有项目用户发送这类请求，查找自己提供的api
		if e := mysql.DB.Order("updated_at DESC").Offset(m.Offset).Limit(m.Limit).Where("A_name like ? AND A_type = ? AND PU_uid = ?", search, 3, userId).Find(&apiRecords).Error; e != nil {
			response.Fail(c, nil, "查找Api时出错")
			return
		}
	case "":
		//没有指定类型，查找所有api
		if e := mysql.DB.Order("updated_at DESC").Offset(m.Offset).Limit(m.Limit).Where("A_name like ?", search).Find(&apiRecords).Error; e != nil {
			response.Fail(c, nil, "查找Api时出错")
			return
		}
	default:
		response.Fail(c, nil, "请求的api类型未知")
		return
	}
	//生成响应报文
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
