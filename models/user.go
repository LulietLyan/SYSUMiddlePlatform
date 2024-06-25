package models

import "time"

// 用户基本属性
type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt  *time.Time `sql:"index"`
	U_uid           uint   `gorm:"primary_key;column:U_uid" json:"U_uid"`
	U_password      string `gorm:"column:U_password;type:VARCHAR(64)" json:"U_password"`
	U_username      string `gorm:"column:U_username;type:VARCHAR(64)" json:"U_username"`
	U_type          uint   `gorm:"column:U_type" json:"U_type"`
	U_mysqlUserName string `gorm:"column:U_mysqlUserName" json:"U_mysqlUserName"`
	U_mysqlUserPwd  string `gorm:"column:U_mysqlUserPwd" jsom:"U_mysqlUserPwd"`
}

func (User) TableName() string {
	return "User" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
