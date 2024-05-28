package models

import "time"

// 用户基本属性
type ProjectUser struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `sql:"index"`
	U_uid          uint       `gorm:"column:U_uid" json:"U_uid"`
	PU_uid         uint       `gorm:"primary_key;column:PU_uid" json:"PU_uid"`
	PU_description string     `gorm:"column:PU_description;type:VARCHAR(8192)" json:"PU_description"`
	PU_write_url   string     `gorm:"column:PU_write_url;type:VARCHAR(512)" json:"PU_write_url"`
	PU_logo_url    string     `gorm:"column:PU_logo_url;type:VARCHAR(1024)" json:"PU_logo_url"`
	PU_email       string     `gorm:"column:PU_email;type:VARCHAR(128)" json:"PU_email"`
}

func (ProjectUser) TableName() string {
	return "ProjectUser" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
