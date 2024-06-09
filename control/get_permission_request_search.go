package control

import (
	"backend/mysql"
	"backend/response"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func GetPermissionRequestSearch(c *gin.Context) {
	type msg struct {
		Limit     uint `json:"limit"`
		Offset    uint `json:"offset"`
		ProjectId uint `json:"projectname"`
		TableId   uint `json:"tablename"`
	}
	var m msg
	if e := c.ShouldBindQuery(&m); e == nil {
		if m.Limit == 0 {
			m.Limit = 15
		}
		var records []struct {
			Requested string    `gorm:"column:U_username"`
			PT_name   string    `gorm:"column:PT_name"`
			PR_uid    uint      `gorm:"column:PR_uid"`
			PR_level  uint      `gorm:"column:PR_level"`
			Requester string    `gorm:"column:requester_pu_name"`
			Time      time.Time `gorm:"column:updated_at"`
		}
		sql := ` SELECT T1.U_username,T1.PT_name,T1.PR_uid,T1.PR_level,T1.updated_at,User.U_username requester_pu_name FROM
		(SELECT User.U_username,ProjectTable.PT_name,PermissionRequest.PR_uid,PermissionRequest.PR_level,PermissionRequest.PU_uid,PermissionRequest.updated_at
			FROM ProjectTable,ProjectUser,User,PermissionRequest
			WHERE User.U_uid=ProjectUser.U_uid AND ProjectUser.PU_uid=ProjectTable.PU_uid AND PermissionRequest.PT_uid=ProjectTable.PT_uid AND PermissionRequest.PR_status=1 `
		if m.ProjectId != 0 {
			sql += fmt.Sprintf(" AND ProjectUser.PU_uid = %d ", m.ProjectId)
		}
		if m.TableId != 0 {
			sql += fmt.Sprintf(" AND ProjectTable.PT_uid = %d ", m.TableId)
		}
		sql += `ORDER BY PermissionRequest.updated_at DESC
			LIMIT ?
			OFFSET ?
		)T1,User,ProjectUser
	WHERE T1.PU_uid=ProjectUser.PU_uid AND ProjectUser.U_uid=User.U_uid`
		if e := mysql.DB.Raw(sql, m.Limit, m.Offset).Scan(&records).Error; e != nil {
			response.Fail(c, nil, "")
			return
		}
		type message struct {
			Title       string `json:"title"`
			Id          uint   `json:"id"`
			Content     string `json:"content"`
			Projectname string `json:"projectname"`
			Time        string `json:"time"`
		}
		var messages []message
		for _, record := range records {
			var levels = [2]string{"只读", "读写"}
			content := fmt.Sprintf("申请数据表“%s”(%s)的%s权限", record.PT_name, record.Requested, levels[record.PR_level-1])
			messages = append(messages, message{Title: "权限请求", Id: record.PR_uid, Content: content, Projectname: record.Requester, Time: record.Time.Format("2006-01-02 15:04")})
		}
		response.Success(c, gin.H{"messages": messages}, "")
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
