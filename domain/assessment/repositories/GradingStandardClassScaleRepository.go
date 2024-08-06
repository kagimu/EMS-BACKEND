package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
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

	if err := db.Scopes(filters.ByAnyField("class_number", classNumber), filters.ByAnyField("scale_number", scaleNumber)).First(&grades).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return grades, nil
		}

		return grades, err
	}

	return grades, nil
}

func AllGradingStandardClassScale(schoolNumber, gradingStandardNumber, classNumber, scaleNumber string, trn ...*gorm.DB) ([]models.GradingStandardClassScale, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var grades []models.GradingStandardClassScale

	if err := db.
		Joins("join grading_standards on grading_standards.number = grading_standard_class_scales.grading_standard_number").
		Scopes(filters.ByAnyField("grading_standard_number", gradingStandardNumber), filters.ByAnyField("school_number", schoolNumber), filters.ByAnyField("scale_number", scaleNumber), filters.ByAnyField("class_number", classNumber)).
		Find(&grades).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return grades, nil
		}

		return grades, err
	}

	return grades, nil
}
