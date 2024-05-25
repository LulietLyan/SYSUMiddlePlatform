package models

import "time"

// 用户基本属性
type ActivationCode struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	AC_uid    uint       `gorm:"primary_key;column:AC_uid" json:"AC_uid"`
	AC_code   string     `gorm:"column:AC_code;type:VARCHAR(256)" json:"AC_code"`
	AC_usable uint       `gorm:"column:AC_usable" json:"AC_usable"`
	AC_type   uint       `gorm:"column:AC_type" json:"AC_type"`
}

func (ActivationCode) TableName() string {
	return "ActivationCode" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
