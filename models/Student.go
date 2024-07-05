package models

type Student struct {
	StudentId       uint   `gorm:"primary_key;column:student_id" json:"student_id"`
	StudentType     string `gorm:"column:student_type" json:"student_type"`
	Gender          string `gorm:"column:gender" json:"gender"`
	Ethnicity       string `gorm:"column:ethnicity" json:"ethnicity"`
	BirthStr        string `gorm:"column:birth_str" json:"birth_str"`
	EducationLevel  string `gorm:"column:education_level" json:"education_level"`
	PoliticalStatus string `gorm:"column:political_status" json:"political_status"`
	Hometown        string `gorm:"column:hometown" json:"hometown"`
	GaokaoScore     uint   `gorm:"column:gaokao_score" json:"gaokao_score"`
	Grade           uint   `gorm:"column:grade" json:"grade"`
	Class           uint   `gorm:"column:class" json:"class"`
}

func (Student) TableName() string {
	return "Student" // 让gorm使用“Admin”作为表名，而不是“Admins”，避免不必要的麻烦
}
