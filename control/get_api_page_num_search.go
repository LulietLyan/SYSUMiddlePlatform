package control

import (
	"backend/mysql"
	"backend/response"

	"github.com/gin-gonic/gin"
)

func GetApiPageNumSearch(c *gin.Context) {
	type msg struct {
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
	var count int64
	if aType > 0 {
		if e := mysql.DB.Order("updated_at DESC").Limit(m.Limit).Where("A_name like ? AND A_type = ?", search, aType).Count(&count).Error; e != nil {
			response.Fail(c, nil, "查找Api数量时出错")
			return
		}
	} else {
		if e := mysql.DB.Order("updated_at DESC").Limit(m.Limit).Where("A_name like ?", search).Count(&count).Error; e != nil {
			response.Fail(c, nil, "查找Api数量时出错")
			return
		}
	}
	pages := count / int64(m.Limit)
	if count%int64(m.Limit) != 0 {
		pages++
	}
	response.Success(c, gin.H{"pages": pages}, "")
}
