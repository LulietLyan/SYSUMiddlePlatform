package control

import (
	"backend/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type returnAvg struct {
	PublicRequiredGpa            string `json:"public_required_gpa"`
	SpecializedElectiveGpa       string `json:"specialized_elective_gpa"`
	SpecializedRequiredGpa       string `json:"specialized_required_gpa"`
	PartyBuildingAwards          string `json:"party_building_awards"`
	AcademicCompetitions         string `json:"academic_competitions"`
	ArtCompetitions              string `json:"art_competitions"`
	SportsCompetitions           string `json:"sports_competitions"`
	EntrepreneurshipCompetitions string `json:"entrepreneurship_competitions"`
	AcademicAchievements         string `json:"academic_achievements"`
	HighLevelPapers              string `json:"high_level_papers"`
	VolunteerHours               string `json:"volunteer_hours"`
	Patents                      string `json:"patents"`
	SoftwareCopyrights           string `json:"software_copyrights"`
	MonographsPublished          string `json:"monographs_published"`
}

func get_all_avg(student_id int) (returnAvgs []returnAvg, err error) {
	// 数据库连接
	tmp, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			"root", "2024", "47.120.73.205", 3307, "source"))
	if err != nil {
		return nil, err
	}

	SelectSQL1 := fmt.Sprintf("SELECT education_level, grade, class "+
		"FROM StudentAchievements WHERE student_id = %d;", student_id)
	var SQLResult1 []struct {
		Level  string `gorm:"column:education_level"`
		Grade  int    `gorm:"column:grade"`
		Class_ int    `gorm:"column:class"`
	}
	if err := tmp.Raw(SelectSQL1).Scan(&SQLResult1).Error; err != nil {
		return nil, err
	}

	SelectSQL21 := fmt.Sprintf("SELECT AVG(public_required_gpa) as public_required_gpa, "+
		"AVG(specialized_elective_gpa) as specialized_elective_gpa, "+
		"AVG(specialized_required_gpa) as specialized_required_gpa, "+
		"AVG(party_building_awards) as party_building_awards, "+
		"AVG(academic_competitions) as academic_competitions, "+
		"AVG(art_competitions) as art_competitions, "+
		"AVG(sports_competitions) as sports_competitions, "+
		"AVG(entrepreneurship_competitions) as entrepreneurship_competitions, "+
		"AVG(academic_achievements) as academic_achievements, "+
		"AVG(high_level_papers) as high_level_papers, "+
		"AVG(volunteer_hours) as volunteer_hours, "+
		"AVG(patents) as patents, "+
		"AVG(software_copyrights) as software_copyrights, "+
		"AVG(monographs_published) as monographs_published "+
		"FROM StudentAchievements WHERE student_id = %d;", student_id)
	var SQLResult21 []struct {
		PublicRequiredGpa            string `gorm:"column:public_required_gpa"`
		SpecializedElectiveGpa       string `gorm:"column:specialized_elective_gpa"`
		SpecializedRequiredGpa       string `gorm:"column:specialized_required_gpa"`
		PartyBuildingAwards          string `gorm:"column:party_building_awards"`
		AcademicCompetitions         string `gorm:"column:academic_competitions"`
		ArtCompetitions              string `gorm:"column:art_competitions"`
		SportsCompetitions           string `gorm:"column:sports_competitions"`
		EntrepreneurshipCompetitions string `gorm:"column:entrepreneurship_competitions"`
		AcademicAchievements         string `gorm:"column:academic_achievements"`
		HighLevelPapers              string `gorm:"column:high_level_papers"`
		VolunteerHours               string `gorm:"column:volunteer_hours"`
		Patents                      string `gorm:"column:patents"`
		SoftwareCopyrights           string `gorm:"column:software_copyrights"`
		MonographsPublished          string `gorm:"column:monographs_published"`
	}
	if err := tmp.Raw(SelectSQL21).Scan(&SQLResult21).Error; err != nil {
		return nil, err
	}

	SelectSQL22 := fmt.Sprintf("SELECT AVG(public_required_gpa) as public_required_gpa, "+
		"AVG(specialized_elective_gpa) as specialized_elective_gpa, "+
		"AVG(specialized_required_gpa) as specialized_required_gpa, "+
		"AVG(party_building_awards) as party_building_awards, "+
		"AVG(academic_competitions) as academic_competitions, "+
		"AVG(art_competitions) as art_competitions, "+
		"AVG(sports_competitions) as sports_competitions, "+
		"AVG(entrepreneurship_competitions) as entrepreneurship_competitions, "+
		"AVG(academic_achievements) as academic_achievements, "+
		"AVG(high_level_papers) as high_level_papers, "+
		"AVG(volunteer_hours) as volunteer_hours, "+
		"AVG(patents) as patents, "+
		"AVG(software_copyrights) as software_copyrights, "+
		"AVG(monographs_published) as monographs_published "+
		"FROM StudentAchievements WHERE education_level = '%s' and grade = %d and class = %d",
		SQLResult1[0].Level, SQLResult1[0].Grade, SQLResult1[0].Class_)
	var SQLResult22 []struct {
		PublicRequiredGpa            string `gorm:"column:public_required_gpa"`
		SpecializedElectiveGpa       string `gorm:"column:specialized_elective_gpa"`
		SpecializedRequiredGpa       string `gorm:"column:specialized_required_gpa"`
		PartyBuildingAwards          string `gorm:"column:party_building_awards"`
		AcademicCompetitions         string `gorm:"column:academic_competitions"`
		ArtCompetitions              string `gorm:"column:art_competitions"`
		SportsCompetitions           string `gorm:"column:sports_competitions"`
		EntrepreneurshipCompetitions string `gorm:"column:entrepreneurship_competitions"`
		AcademicAchievements         string `gorm:"column:academic_achievements"`
		HighLevelPapers              string `gorm:"column:high_level_papers"`
		VolunteerHours               string `gorm:"column:volunteer_hours"`
		Patents                      string `gorm:"column:patents"`
		SoftwareCopyrights           string `gorm:"column:software_copyrights"`
		MonographsPublished          string `gorm:"column:monographs_published"`
	}
	fmt.Println(SelectSQL22)
	if err := tmp.Raw(SelectSQL22).Scan(&SQLResult22).Error; err != nil {
		fmt.Printf("Error executing SQL: %v\n", err)
		return nil, err
	}

	SelectSQL23 := fmt.Sprintf("SELECT AVG(public_required_gpa) as public_required_gpa, "+
		"AVG(specialized_elective_gpa) as specialized_elective_gpa, "+
		"AVG(specialized_required_gpa) as specialized_required_gpa, "+
		"AVG(party_building_awards) as party_building_awards, "+
		"AVG(academic_competitions) as academic_competitions, "+
		"AVG(art_competitions) as art_competitions, "+
		"AVG(sports_competitions) as sports_competitions, "+
		"AVG(entrepreneurship_competitions) as entrepreneurship_competitions, "+
		"AVG(academic_achievements) as academic_achievements, "+
		"AVG(high_level_papers) as high_level_papers, "+
		"AVG(volunteer_hours) as volunteer_hours, "+
		"AVG(patents) as patents, "+
		"AVG(software_copyrights) as software_copyrights, "+
		"AVG(monographs_published) as monographs_published "+
		"FROM StudentAchievements WHERE education_level = '%s' and grade = %d;",
		SQLResult1[0].Level, SQLResult1[0].Grade)
	var SQLResult23 []struct {
		PublicRequiredGpa            string `gorm:"column:public_required_gpa"`
		SpecializedElectiveGpa       string `gorm:"column:specialized_elective_gpa"`
		SpecializedRequiredGpa       string `gorm:"column:specialized_required_gpa"`
		PartyBuildingAwards          string `gorm:"column:party_building_awards"`
		AcademicCompetitions         string `gorm:"column:academic_competitions"`
		ArtCompetitions              string `gorm:"column:art_competitions"`
		SportsCompetitions           string `gorm:"column:sports_competitions"`
		EntrepreneurshipCompetitions string `gorm:"column:entrepreneurship_competitions"`
		AcademicAchievements         string `gorm:"column:academic_achievements"`
		HighLevelPapers              string `gorm:"column:high_level_papers"`
		VolunteerHours               string `gorm:"column:volunteer_hours"`
		Patents                      string `gorm:"column:patents"`
		SoftwareCopyrights           string `gorm:"column:software_copyrights"`
		MonographsPublished          string `gorm:"column:monographs_published"`
	}
	if err := tmp.Raw(SelectSQL23).Scan(&SQLResult23).Error; err != nil {
		return nil, err
	}
	if err := tmp.Close(); err != nil {
		return nil, err
	}

	returnAvgs = append(returnAvgs, returnAvg{
		SQLResult21[0].PublicRequiredGpa,
		SQLResult21[0].SpecializedElectiveGpa,
		SQLResult21[0].SpecializedRequiredGpa,
		SQLResult21[0].PartyBuildingAwards,
		SQLResult21[0].AcademicCompetitions,
		SQLResult21[0].ArtCompetitions,
		SQLResult21[0].SportsCompetitions,
		SQLResult21[0].EntrepreneurshipCompetitions,
		SQLResult21[0].AcademicAchievements,
		SQLResult21[0].HighLevelPapers,
		SQLResult21[0].VolunteerHours,
		SQLResult21[0].Patents,
		SQLResult21[0].SoftwareCopyrights,
		SQLResult21[0].MonographsPublished})
	returnAvgs = append(returnAvgs, returnAvg{
		SQLResult22[0].PublicRequiredGpa,
		SQLResult22[0].SpecializedElectiveGpa,
		SQLResult22[0].SpecializedRequiredGpa,
		SQLResult22[0].PartyBuildingAwards,
		SQLResult22[0].AcademicCompetitions,
		SQLResult22[0].ArtCompetitions,
		SQLResult22[0].SportsCompetitions,
		SQLResult22[0].EntrepreneurshipCompetitions,
		SQLResult22[0].AcademicAchievements,
		SQLResult22[0].HighLevelPapers,
		SQLResult22[0].VolunteerHours,
		SQLResult22[0].Patents,
		SQLResult22[0].SoftwareCopyrights,
		SQLResult22[0].MonographsPublished})
	returnAvgs = append(returnAvgs, returnAvg{
		SQLResult23[0].PublicRequiredGpa,
		SQLResult23[0].SpecializedElectiveGpa,
		SQLResult23[0].SpecializedRequiredGpa,
		SQLResult23[0].PartyBuildingAwards,
		SQLResult23[0].AcademicCompetitions,
		SQLResult23[0].ArtCompetitions,
		SQLResult23[0].SportsCompetitions,
		SQLResult23[0].EntrepreneurshipCompetitions,
		SQLResult23[0].AcademicAchievements,
		SQLResult23[0].HighLevelPapers,
		SQLResult23[0].VolunteerHours,
		SQLResult23[0].Patents,
		SQLResult23[0].SoftwareCopyrights,
		SQLResult23[0].MonographsPublished})
	return returnAvgs, nil
}

func GetAllAvg(c *gin.Context) {
	type msg struct {
		Id int `json:"id"`
	}

	var m msg
	if e := c.ShouldBindJSON(&m); e == nil {
		if returnAvgs, err := get_all_avg(m.Id); err == nil {
			response.Success(c, gin.H{"avg1": returnAvgs[0], "avg2": returnAvgs[1], "avg3": returnAvgs[2]}, "")
		}
		response.Fail(c, nil, "")
	} else { //JSON解析失败
		response.Fail(c, nil, "数据格式错误!")
	}
}
