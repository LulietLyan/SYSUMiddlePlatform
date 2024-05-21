package models

import "time"

// 用户基本属性
type AnalyticalUser struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	AU_uid      uint       `gorm:"primary_key;column:AU_uid" json:"AU_uid"`
	AU_password string     `gorm:"column:AU_password;type:VARCHAR(64)" json:"AU_password"`
	AU_username string     `gorm:"column:AU_username;type:VARCHAR(64)" json:"AU_username"`
	AU_phone    string     `gorm:"column:AU_phone;type:VARCHAR(20)" json:"AU_phone"`
	AU_std_uid  string     `gorm:"column:AU_std_uid;type:VARCHAR(20)" json:"AU_std_uid"`
	AU_email    string     `gorm:"column:AU_email;type:VARCHAR(128)" json:"AU_email"`
	AU_realname string     `gorm:"column:AU_realname;type:VARCHAR(64)" json:"AU_realname"`
}

func (AnalyticalUser) TableName() string {
	return "AnalyticalUser" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
