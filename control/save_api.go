package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"
	"fmt"

	"github.com/gin-gonic/gin"
)

func SaveApi(c *gin.Context) {
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
	//解析参数
	type msg struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Url         string `json:"url"`
		Desc        string `json:"desc"`
		Request     string `json:"request"`
		Response    string `json:"response"`
		Id          int64  `json:"id"`
		Projectname string `json:"projectname"`
	}
	var m msg
	if e := c.ShouldBindJSON(&m); e != nil {
		response.Fail(c, nil, "数据格式错误!")
		return
	}
	var aType uint
	fmt.Println(m.Type)
	switch m.Type {
	case "Midtable":
		aType = 1
	case "Require":
		aType = 2
	case "User":
		aType = 3
	case "Me":
		aType = 3
	default:
		response.Fail(c, nil, "未知的Api类型!")
		return
	}
	if m.Id < 0 {
		//Id小于0表示新建
		aRecord := models.Api{A_url: m.Url, A_parameter: m.Request, A_respond: m.Response, A_description: m.Desc, A_type: aType, A_name: m.Name}
		if aType == 3 {
			var puRecord models.ProjectUser
			if e := mysql.DB.Where("U_uid = ?", userId).First(&puRecord).Error; e != nil {
				response.Fail(c, nil, "查找项目用户时失败")
				return
			}
			aRecord.PU_uid = puRecord.PU_uid
		}
		if e := mysql.DB.Create(&aRecord).Error; e != nil {
			response.Fail(c, nil, "插入新Api信息时出错")
			return
		}
		response.Success(c, nil, "")
		return
	} else {
		//Id大于0表示修改
		if identity == "Developer" {
			if aType == 3 { //项目用户请求更新自己的api
				//检查是否有权限
				var uRecord models.User
				if e := mysql.DB.Where("U_uid = ?", userId).First(&uRecord).Error; e == nil {
					if uRecord.U_username != m.Projectname {
						response.Fail(c, nil, "不能修改其他项目的表")
						return
					}
				} else {
					response.Fail(c, nil, "项目名不存在")
					return
				}
				//保存新值
				aRecord := models.Api{A_uid: uint(m.Id), A_url: m.Url, A_parameter: m.Request, A_respond: m.Response, A_description: m.Desc, A_type: aType, A_name: m.Name, PU_uid: userId}
				if e := mysql.DB.Save(&aRecord).Error; e != nil {
					response.Fail(c, nil, "更新Api时出错")
					return
				}
				response.Success(c, nil, "")
				return
			} else if aType == 2 { //项目用户上传url
				var puRecord models.ProjectUser
				if e := mysql.DB.Raw("SELECT ProjectUser.* FROM User,ProjectUser WHERE User.U_uid=ProjectUser.U_uid and User.U_uid = ? and User.U_username = ?", userId, m.Projectname).First(&puRecord).Error; e != nil {
					response.Fail(c, nil, "查找项目用户记录时出错")
					return
				}
				puRecord.PU_write_url = m.Url
				if e := mysql.DB.Save(&puRecord).Error; e != nil {
					response.Fail(c, nil, "更新项目写回url路径时出错")
					return
				}
				response.Success(c, nil, "")
				return
			} else {
				response.Fail(c, nil, "项目端用户不允许修改此类Api")
				return
			}
		} else if identity == "Admin" {
			if aType == 3 {
				response.Fail(c, nil, "管理员不允许修改此类Api")
				return
			} else {
				aRecord := models.Api{A_uid: uint(m.Id), A_url: m.Url, A_parameter: m.Request, A_respond: m.Response, A_description: m.Desc, A_type: aType, A_name: m.Name}
				if e := mysql.DB.Save(&aRecord).Error; e != nil {
					response.Fail(c, nil, "更新Api时出错")
					return
				}
				response.Success(c, nil, "")
				return
			}
		} else {
			response.Fail(c, nil, "未定义行为")
			return
		}
	}

}
