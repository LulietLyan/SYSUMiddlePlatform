package models

import "time"

// 用户基本属性
type Permission struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt *time.Time `sql:"index"`
	P_uid   uint `gorm:"column:P_uid;primary_key" json:"P_uid"`
	PU_uid  uint `gorm:"column:PU_uid" json:"PU_uid"`
	PT_uid  uint `gorm:"column:PT_uid" json:"PT_uid"`
	P_level uint `gorm:"column:P_level" json:"P_level"`
}

func (Permission) TableName() string {
	return "Permission" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
