package models

import "time"

// 用户基本属性
type DingdingProjectUser struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	DPU_uid   uint       `gorm:"column:DPU_uid;primary_key" json:"DPU_uid"`
	DPU_phone string     `gorm:"column:DPU_phone;type:VARCHAR(20)" json:"DPU_phone"`
	PU_uid    uint       `gorm:"column:PU_uid" json:"PU_uid"`
	PM_uid    uint       `gorm:"column:PM_uid" json:"PM_uid"`
}

func (DingdingProjectUser) TableName() string {
	return "DingdingProjectUser" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
