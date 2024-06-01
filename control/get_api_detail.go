package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetApiDetail(c *gin.Context) {
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
		Id uint `json:"id"`
	}

	var m msg
	if e := c.ShouldBindQuery(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}
	//查找apiRecord
	var apiRecord models.Api
	if e := mysql.DB.Where("A_uid = ?", m.Id).First(&apiRecord).Error; e != nil {
		response.Fail(c, nil, "Api不存在")
		return
	}
	//生成响应
	type apiInfo struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Url         string `json:"url"`
		Desc        string `json:"desc"`
		Request     string `json:"request"`
		Response    string `json:"response"`
		Id          uint   `json:"id"`
		Projectname string `json:"projectname"`
	}
	var apiTypes = [4]string{"Midtable", "Require", "User", "Me"}
	typeindex := apiRecord.A_type - 1 //1 2 3 --> 0 1 2
	if identity == "Developer" && apiRecord.A_type == 3 && apiRecord.PU_uid == userId {
		typeindex = 3 //标记为自己提供的api
	}
	var queryProjectName struct {
		ProjectName string `gorm:"column:U_username"`
	}
	if apiRecord.A_type == 3 {
		//如果是用户提供的api，查找用户名并添加到响应中
		if err := mysql.DB.Raw("SELECT User.U_username FROM ProjectUser,User Where ProjectUser.U_uid = User.U_uid and ProjectUser.PU_uid = ?", apiRecord.PU_uid).First(&queryProjectName).Error; err != nil {
			response.Fail(c, nil, "查询提供api的项目名时出错!")
			return
		}
	}
	data := apiInfo{Name: apiRecord.A_name, Type: apiTypes[typeindex], Url: apiRecord.A_url, Desc: apiRecord.A_description, Request: apiRecord.A_parameter, Response: apiRecord.A_respond, Id: apiRecord.A_uid, Projectname: queryProjectName.ProjectName}
	response.Success(c, gin.H{"apiInfo": data}, "")
}
