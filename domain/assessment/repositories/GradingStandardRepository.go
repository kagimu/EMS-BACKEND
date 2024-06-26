package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
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

func FindGradingStandardByNameAndClassNumber(name string, campusClassNumber string) (models.GradingStandard, error) {
	db := config.DB
	var gradingStandard models.GradingStandard

	if err := db.Where("name = ? AND campus_class_number = ?", name, campusClassNumber).First(&gradingStandard).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gradingStandard, nil
		}

		return gradingStandard, err
	}

	return gradingStandard, nil
}

func FindGradingStandardByCodeAndClassNumber(code int64, campusClassNumber string) (models.GradingStandard, error) {
	db := config.DB
	var gradingStandard models.GradingStandard

	if err := db.Where("code = ? AND campus_class_number = ?", code, campusClassNumber).First(&gradingStandard).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gradingStandard, nil
		}

		return gradingStandard, err
	}

	return gradingStandard, nil
}

func GetGradingStandardByNumber(gradingStandardNumber string) (models.GradingStandard, error) {
	db := config.DB
	var gradingStandard models.GradingStandard
	if err := db.Where("number = ?", gradingStandardNumber).First(&gradingStandard).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gradingStandard, nil
		}
		return gradingStandard, err
	}
	return gradingStandard, nil
}
