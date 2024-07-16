package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	"gorm.io/gorm"
)

func SaveAssessmentGradeScale(assessmentGradeScale models.AssessmentGradeScale, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&assessmentGradeScale).Error
}

func GetAssessmentGradeScaleByCodeAndGradeNumber(code int64, assessmentGradeNumber string, trn ...*gorm.DB) (models.AssessmentGradeScale, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentGrade models.AssessmentGradeScale

	if err := db.Where("code = ? AND assessment_grade_number = ?", code, assessmentGradeNumber).First(&assessmentGrade).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentGrade, nil
		}

		return assessmentGrade, err
	}

	return assessmentGrade, nil
}

func FindAssessmentGradeScaleByGradeNumber(assessmentGradeNumber string, trn ...*gorm.DB) ([]models.AssessmentGradeScale, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentGrades []models.AssessmentGradeScale

	if err := db.Where("assessment_grade_number = ?", assessmentGradeNumber).Find(&assessmentGrades).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentGrades, nil
		}

		return assessmentGrades, err
	}

	return assessmentGrades, nil
}

func RemoveAssessmentGradeScale(assessmentGradeScale models.AssessmentGradeScale, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Delete(&assessmentGradeScale).Error
}

func GetAssessmentGradeScaleByNameAndGradeNumber(name, assessmentGradeNumber string, trn ...*gorm.DB) (models.AssessmentGradeScale, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentGrade models.AssessmentGradeScale

	if err := db.Where("name = ? AND assessment_grade_number = ?", name, assessmentGradeNumber).First(&assessmentGrade).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentGrade, nil
		}

		return assessmentGrade, err
	}

	return assessmentGrade, nil
}
