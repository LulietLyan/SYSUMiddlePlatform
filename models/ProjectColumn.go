package models

import "time"

// 用户基本属性
type ProjectColumn struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `sql:"index"`
	PC_uid         uint       `gorm:"column:PC_uid;primary_key" json:"PC_uid"`
	PC_name        string     `gorm:"column:PC_name;type:VARCHAR(64)" json:"PC_name"`
	PC_description string     `gorm:"column:PC_description;type:VARCHAR(1024)" json:"PC_description"`
	PT_uid         uint       `gorm:"column:PT_uid" json:"PT_uid"`
}

func (ProjectColumn) TableName() string {
	return "ProjectColumn" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
