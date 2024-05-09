package models

import "time"

// 用户基本属性
type Admin struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `sql:"index"`
	Admin_uid      uint       `gorm:"primary_key;column:Admin_uid" json:"Admin_uid"`
	Admin_password string     `gorm:"column:Admin_password;type:VARCHAR(64)" json:"Admin_password"`
	Admin_username string     `gorm:"column:Admin_username;type:VARCHAR(64)" json:"Admin_username"`
}

func (Admin) TableName() string {
	return "Admin" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
