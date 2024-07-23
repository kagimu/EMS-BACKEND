package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func SaveAssessmentGrade(assessmentGrade models.AssessmentGrade, trn ...*gorm.DB) error {
	db := config.DB
	if trn != nil && len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&assessmentGrade).Error
}

func FindAssessmentGradeByClassNumber(classNumber, status string, trn ...*gorm.DB) ([]models.AssessmentGrade, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentGrades []models.AssessmentGrade

	if err := db.
		Joins("join assessment_grade_classes on assessment_grade_classes.assessment_grade_number = assessment_grades.number").
		Scopes(filters.ByAnyField("class_number", classNumber), filters.ByStatus(status)).
		Distinct("assessment_grades.number").
		Find(&assessmentGrades).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentGrades, nil
		}

		return assessmentGrades, err
	}

	return assessmentGrades, nil
}

func GetAssessmentGradeByClassNumberAndStatus(classNumber, status string, trn ...*gorm.DB) (models.AssessmentGrade, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentGrade models.AssessmentGrade

	if err := db.
		Joins("join assessment_grade_classes on assessment_grade_classes.assessment_grade_number = assessment_grades.number").
		Scopes(filters.ByAnyField("assessment_grade_classes.class_number", classNumber), filters.ByAnyField("assessment_grades.status", status)).
		Distinct("assessment_grades.number", "assessment_grades.status", "assessment_grades.assessment_area_count", "assessment_grades.user_number", "assessment_grades.school_number", "assessment_grades.created_at", "assessment_grades.updated_at").
		First(&assessmentGrade).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentGrade, nil
		}

		return assessmentGrade, err
	}

	return assessmentGrade, nil
}

func FindAssessmentGradeByNumber(assessmentGradeNumber string, trn ...*gorm.DB) (models.AssessmentGrade, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentGrade models.AssessmentGrade

	if err := db.Where("number = ? ", assessmentGradeNumber).First(&assessmentGrade).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentGrade, nil
		}

		return assessmentGrade, err
	}

	return assessmentGrade, nil
}

func RemoveAssessmentGrade(assessmentGrade models.AssessmentGrade, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Delete(&assessmentGrade).Error
}
