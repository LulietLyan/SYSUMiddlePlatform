package models

import "time"

// 用户基本属性
type PermissionRequest struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt *time.Time `sql:"index"`
	PR_uid    uint `gorm:"column:PR_uid;primary_key" json:"PR_uid"`
	PU_uid    uint `gorm:"column:PU_uid" json:"PU_uid"`
	PT_uid    uint `gorm:"column:PT_uid" json:"PT_uid"`
	PR_level  uint `gorm:"column:PR_level" json:"PR_level"`
	PR_status uint `gorm:"column:PR_status" json:"PR_status"`
}

func (PermissionRequest) TableName() string {
	return "PermissionRequest" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
