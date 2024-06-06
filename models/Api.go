package models

import "time"

// 用户基本属性
type Api struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt     *time.Time `sql:"index"`
	A_uid         uint   `gorm:"column:A_uid;primary_key" json:"A_uid"`
	A_url         string `gorm:"column:A_url;type:VARCHAR(1024)" json:"A_url"`
	A_parameter   string `gorm:"column:A_parameter;type:VARCHAR(5000)" json:"A_parameter"`
	A_respond     string `gorm:"column:A_respond;type:VARCHAR(5000)" json:"A_respond"`
	A_description string `gorm:"column:A_description;type:VARCHAR(2048)" json:"A_description"`
	A_type        uint   `gorm:"column:A_type" json:"A_type"`
	A_name        string `gorm:"column:A_name;type:VARCHAR(64)" json:"A_name"`
	PU_uid        uint   `gorm:"column:PU_uid" json:"PU_uid"`
}

func (Api) TableName() string {
	return "Api" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
