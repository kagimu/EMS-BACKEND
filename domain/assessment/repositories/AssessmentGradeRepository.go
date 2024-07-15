package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func GetAssessmentGradeByNumber(assessmentGradeNumber string, trn ...*gorm.DB) (models.AssessmentGrade, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentGrade models.AssessmentGrade

	if err := db.Where("number = ?", assessmentGradeNumber).First(&assessmentGrade).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentGrade, nil
		}

		return assessmentGrade, err
	}

	return assessmentGrade, nil
}

func SaveAssessmentGrade(assessmentGrade models.AssessmentGrade, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&assessmentGrade).Error
}

func FindAssessmentGradeByClassNumber(classNumber string, trn ...*gorm.DB) ([]models.AssessmentGrade, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentGrades []models.AssessmentGrade

	if err := db.Scopes(filters.ByAnyField("class_number", classNumber)).Find(&assessmentGrades).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentGrades, nil
		}

		return assessmentGrades, err
	}

	return assessmentGrades, nil
}
