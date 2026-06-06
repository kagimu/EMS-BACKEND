package main

import (
	"fmt"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
)

func main() {
	config.ConnectDB() // Assume we can connect or we might need to load env

	db := config.DB
	if db == nil {
		fmt.Println("DB is nil")
		return
	}

	var enrollmentInfos []struct {
		UserNumber   string
		Adn          string
		Class        string
		Level        string
		AcademicYear string
		StudyPeriod  string
	}

	err := db.Table("students").
		Select("students.user_number as user_number, students.adn as adn, classes.name as class, levels.name as level, academic_years.name as academic_year, study_period_names.name as study_period").
		Joins("JOIN enrollments ON enrollments.student_number = students.number").
		Joins("JOIN campus_classes ON campus_classes.number = enrollments.campus_class_number").
		Joins("JOIN curriculum_classes ON curriculum_classes.number = campus_classes.curriculum_class_number").
		Joins("JOIN classes ON classes.number = curriculum_classes.class_number").
		Joins("JOIN levels ON levels.number = classes.level_number").
		Joins("JOIN study_periods ON study_periods.number = enrollments.study_period_number").
		Joins("JOIN study_period_names ON study_period_names.number = study_periods.study_period_name_number").
		Joins("JOIN academic_years ON academic_years.number = study_periods.academic_year_number").
		Where("enrollments.status = ?", constants.StatusActive).
		Limit(5).
		Scan(&enrollmentInfos).Error

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Fetched %d enrollment infos\n", len(enrollmentInfos))
	for _, info := range enrollmentInfos {
		fmt.Printf("%+v\n", info)
	}
}
