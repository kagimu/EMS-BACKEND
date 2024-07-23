package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func SaveGradingStandard(gradingStandard models.GradingStandard, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Save(&gradingStandard).Error; err != nil {
		return err
	}

	return nil

}

func FindGradingStandardByNameAndSchoolNumber(name, schoolNumber string, trn ...*gorm.DB) (models.GradingStandard, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var gradingStandard models.GradingStandard

	if err := db.Where("name = ? AND school_number = ?", name, schoolNumber).First(&gradingStandard).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gradingStandard, nil
		}

		return gradingStandard, err
	}

	return gradingStandard, nil
}

func FindGradingStandardByCodeAndSchoolNumber(code int64, schoolNumber string, trn ...*gorm.DB) (models.GradingStandard, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var gradingStandard models.GradingStandard

	if err := db.Where("code = ? AND school_number = ?", code, schoolNumber).First(&gradingStandard).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gradingStandard, nil
		}

		return gradingStandard, err
	}

	return gradingStandard, nil
}

func GetGradingStandardByNumber(gradingStandardNumber string, trn ...*gorm.DB) (models.GradingStandard, error) {
	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	var gradingStandard models.GradingStandard
	if err := db.Where("number = ?", gradingStandardNumber).First(&gradingStandard).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gradingStandard, nil
		}
		return gradingStandard, err
	}
	return gradingStandard, nil
}

func FindGradingStandardByCampusNumberAndCampusClassNumber(schoolNumber, classNumber, scaleNumber string, trn ...*gorm.DB) ([]models.GradingStandard, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var gradingStandards []models.GradingStandard

	if err := db.
		Joins("join grading_standard_class_scales on grading_standard_class_scales.grading_standard_number = grading_standards.number").
		Scopes(filters.BySchoolNumber(schoolNumber), filters.ByAnyField("grading_standard_class_scales.class_number", classNumber), filters.ByAnyField("grading_standard_class_scales.scale_number", scaleNumber)).
		Distinct("grading_standards.number").
		Find(&gradingStandards).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gradingStandards, nil
		}
		return gradingStandards, err
	}
	return gradingStandards, nil
}

func RemoveGradingStandard(gradingStandard models.GradingStandard, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Unscoped().Delete(&gradingStandard).Error; err != nil {
		return err
	}

	return nil

}
