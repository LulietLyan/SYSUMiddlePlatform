package models

import "time"

// 用户基本属性
type Notifications struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	N_uid     uint       `gorm:"column:N_uid;primary_key" json:"N_uid"`
	N_type    uint       `gorm:"column:N_type" json:"N_type"`
	PU_uid    uint       `gorm:"column:PU_uid" json:"PU_uid"`
	AU_uid    uint       `gorm:"column:AU_uid" json:"AU_uid"`
	N_Title   string     `gorm:"column:N_Title;type:VARCHAR(256)" json:"N_Title"`
	N_Body    string     `gorm:"column:N_Body;type:VARCHAR(8192)" json:"N_Body"`
}

func (Notifications) TableName() string {
	return "Notifications" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
