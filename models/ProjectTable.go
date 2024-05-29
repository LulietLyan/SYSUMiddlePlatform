package models

import "time"

// 用户基本属性
type ProjectTable struct {
	PT_uid               uint `gorm:"column:PT_uid;primary_key" json:"PT_uid"`
	CreatedAt            time.Time
	PT_name              string `gorm:"column:PT_name;type:VARCHAR(64)" json:"PT_name"`
	PT_description       string `gorm:"column:PT_description;type:VARCHAR(8192)" json:"PT_description"`
	PU_uid               uint   `gorm:"column:PU_uid" json:"PU_uid"`
	PT_remote_db_name    string `gorm:"column:PT_remote_db_name;type:VARCHAR(64)" json:"PT_remote_db_name"`
	PT_remote_table_name string `gorm:"column:PT_remote_table_name;type:VARCHAR(64)" json:"PT_remote_table_name"`
	PT_remote_hostname   string `gorm:"column:PT_remote_hostname;type:VARCHAR(64)" json:"PT_remote_hostname"`
	PT_remote_username   string `gorm:"column:PT_remote_username;type:VARCHAR(64)" json:"PT_remote_username"`
	PT_remote_password   string `gorm:"column:PT_remote_password;type:VARCHAR(64)" json:"PT_remote_password"`
	PT_remote_port       uint   `gorm:"column:PT_remote_port" json:"PT_remote_port"`
}

func (ProjectTable) TableName() string {
	return "ProjectTable" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
