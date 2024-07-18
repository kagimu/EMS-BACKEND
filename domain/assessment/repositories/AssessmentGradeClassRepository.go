package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"gorm.io/gorm"
)

func GetAssessmentGradeClassByGradeNumberAndClassNumber(assessmentGradeNumber, classNumber string, trn ...*gorm.DB) (models.AssessmentGradeClass, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentGradeClass models.AssessmentGradeClass

	if err := db.Where("class_number = ? AND assessment_grade_number = ?", classNumber, assessmentGradeNumber).First(&assessmentGradeClass).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentGradeClass, nil
		}

		return assessmentGradeClass, err
	}

	return assessmentGradeClass, nil
}

func SaveAssessmentGradeClass(assessmentGradeClass models.AssessmentGradeClass, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&assessmentGradeClass).Error
}

func FindAssessmentGradeClassByGradeNumber(assessmentGradeNumber string, trn ...*gorm.DB) ([]models.AssessmentGradeClass, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentGradeClasses []models.AssessmentGradeClass

	if err := db.Where("assessment_grade_number = ?", assessmentGradeNumber).Find(&assessmentGradeClasses).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentGradeClasses, nil
		}

		return assessmentGradeClasses, err
	}

	return assessmentGradeClasses, nil
}

func RemoveAssessmentGradeClass(assessmentGradeClass models.AssessmentGradeClass, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Delete(&assessmentGradeClass).Error
}
