package models

import "time"

// 用户基本属性
type DingdingAnalyticalUser struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt *time.Time `sql:"index"`
	DAU_uid   uint   `gorm:"column:DAU_uid;primary_key" json:"DAU_uid"`
	DAU_phone string `gorm:"column:DAU_phone;type:VARCHAR(20)" json:"DAU_phone"`
	AU_uid    uint   `gorm:"column:AU_uid" json:"AU_uid"`
}

func (DingdingAnalyticalUser) TableName() string {
	return "DingdingAnalyticalUser" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
