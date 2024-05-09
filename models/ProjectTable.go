package models

import "time"

// 用户基本属性
type ProjectTable struct {
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            *time.Time `sql:"index"`
	PT_uid               uint       `gorm:"column:PT_uid;primary_key" json:"PT_uid"`
	PT_name              string     `gorm:"column:PT_name;type:VARCHAR(64)" json:"PT_name"`
	PT_description       string     `gorm:"column:PT_description;type:VARCHAR(8192)" json:"PT_description"`
	PT_config            string     `gorm:"column:PT_config;type:TEXT" json:"PT_config"`
	PT_remote_db_name    string     `gorm:"column:PT_remote_db_name;type:VARCHAR(64)" json:"PT_remote_db_name"`
	PT_remote_table_name string     `gorm:"column:PT_remote_table_name;type:VARCHAR(64)" json:"PT_remote_table_name"`
	PU_uid               uint       `gorm:"column:PU_uid" json:"PU_uid"`
}

func (ProjectTable) TableName() string {
	return "ProjectTable" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
