package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetMessagePageNumSearch(c *gin.Context) {
	type msg struct {
		Limit  uint   `json:"limit"`
		Search string `json:"search"`
	}

	var m msg
	if e := c.ShouldBindJSON(&m); e == nil {
		if m.Limit == 0 {
			m.Limit = 15
		}

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
		search := "%" + m.Search + "%"
		var count int64
		switch identity {
		case "Admin":
			if e := mysql.DB.Model(&models.Notifications{}).Where("N_Title like ? OR N_Body like ?", search, search).Count(&count).Error; e != nil {
				response.Fail(c, nil, "查找通知时出错")
				return
			}
		case "Analyzer":
			if e := mysql.DB.Model(&models.Notifications{}).Where("(N_Title like ? OR N_Body like ?) AND (N_type = 2 OR N_type = 3 OR (N_type = 5 AND AU_uid = ?))", search, search, userId).Count(&count).Error; e != nil {
				response.Fail(c, nil, "查找通知时出错")
				return
			}
		case "Developer":
			if e := mysql.DB.Model(&models.Notifications{}).Where("(N_Title like ? OR N_Body like ?) AND (N_type = 1 OR N_type = 3 OR (N_type = 4 AND PU_uid = ?))", search, search, userId).Count(&count).Error; e != nil {
				response.Fail(c, nil, "查找通知时出错")
				return
			}
		default:
			response.Fail(c, nil, "Identity参数为未知值")
			return
		}
		print(count)
		pages := count / int64(m.Limit)
		if count%int64(m.Limit) != 0 {
			pages++
		}
		response.Success(c, gin.H{"pages": pages}, "")

	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
