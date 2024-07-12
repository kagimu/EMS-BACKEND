package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	"gorm.io/gorm"
)

func SaveGradingStandardClassScale(gradingStandardClassScale models.GradingStandardClassScale, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&gradingStandardClassScale).Error
}

func GetGradingStandardClassScaleByClassNumberAndStandard(classNumber, gradingStandardNumber string, trn ...*gorm.DB) (models.GradingStandardClassScale, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var grades models.GradingStandardClassScale

	if err := db.Where("class_number = ? AND grading_standard_number = ?", classNumber, gradingStandardNumber).Find(&grades).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return grades, nil
		}

		return grades, err
	}

	return grades, nil
}

func RemoveGradingStandardClassScale(grades models.GradingStandardClassScale, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Unscoped().Delete(&grades).Error
}

func FindGradingStandardClassScaleByGradingStandardNumber(gradingStandardNumber string, trn ...*gorm.DB) ([]models.GradingStandardClassScale, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var grades []models.GradingStandardClassScale

	if err := db.Where("grading_standard_number = ?", gradingStandardNumber).Find(&grades).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return grades, nil
		}

		return grades, err
	}

	return grades, nil
}

func FindGradingStandardClassScaleByScaleNumber(scaleNumber string, trn ...*gorm.DB) ([]models.GradingStandardClassScale, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var grades []models.GradingStandardClassScale

	if err := db.Where("scale_number = ?", scaleNumber).Find(&grades).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return grades, nil
		}

		return grades, err
	}

	return grades, nil
}

func GetGradingStandardClassScaleByClassNumberAndScaleNumber(classNumber, scaleNumber string, trn ...*gorm.DB) (models.GradingStandardClassScale, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var grades models.GradingStandardClassScale

	if err := db.Where("class_number = ? AND scale_number = ?", classNumber, scaleNumber).Find(&grades).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return grades, nil
		}

		return grades, err
	}

	return grades, nil
}
