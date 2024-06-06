package models

import "time"

// 用户基本属性
type ProjectMember struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt   *time.Time `sql:"index"`
	PM_uid      uint   `gorm:"column:PM_uid;primary_key" json:"PM_uid"`
	PM_name     string `gorm:"column:PM_name;type:VARCHAR(64)" json:"PM_name"`
	PM_phone    string `gorm:"column:PM_phone;type:VARCHAR(20)" json:"PM_phone"`
	PM_email    string `gorm:"column:PM_email;type:VARCHAR(128)" json:"PM_email"`
	PU_uid      uint   `gorm:"column:PU_uid" json:"PU_uid"`
	PM_position string `gorm:"column:PM_position;type:VARCHAR(64)" json:"PM_position"`
}

func (ProjectMember) TableName() string {
	return "ProjectMember" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
