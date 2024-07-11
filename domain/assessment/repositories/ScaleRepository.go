package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	"gorm.io/gorm"
)

func SaveScale(scale models.Scale, trn ...*gorm.DB) (models.Scale, error) {
	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	return scale, db.Save(&scale).Error
}

func GetScaleByNumber(scaleNumber string, trn ...*gorm.DB) (models.Scale, error) {
	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	var scale models.Scale

	if err := db.Where("number = ?", scaleNumber).Find(&scale).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return scale, nil
		}

		return scale, err
	}

	return scale, nil
}

func GetScaleByNameAndSchoolNumber(scaleNumber, schoolNumber string, trn ...*gorm.DB) (models.Scale, error) {
	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	var scale models.Scale

	if err := db.Where("name = ? AND school_number = ? ", scaleNumber, schoolNumber).Find(&scale).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return scale, nil
		}

		return scale, err
	}

	return scale, nil
}

func FindAllScales(schoolNumber string, trn ...*gorm.DB) ([]models.Scale, error) {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	var scales []models.Scale

	if err := db.Where("school_number = ? ", schoolNumber).Find(&scales).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return scales, nil
		}

		return scales, err
	}

	return scales, nil
}

func RemoveScale(scale models.Scale, trn ...*gorm.DB) error {
	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Delete(&scale).Error
}
