package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

func SaveModeOfStudy(modeOfStudy models.ModeOfStudy, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&modeOfStudy).Error
}

func FindModeOfStudy(schoolNumber string, trn ...*gorm.DB) ([]models.ModeOfStudy, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var modeOfStudies []models.ModeOfStudy

	if err := db.Where("school_number = ?", schoolNumber).Find(&modeOfStudies).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return modeOfStudies, nil
		}

		return modeOfStudies, err
	}

	return modeOfStudies, nil
}

func RemoveModeOfStudy(modeOfStudy models.ModeOfStudy, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Delete(&modeOfStudy).Error
}

func GetModeOfStudyByName(name, schoolNumber string, trn ...*gorm.DB) (models.ModeOfStudy, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var modeOfStudy models.ModeOfStudy

	if err := db.Where("school_number = ? AND name = ?", schoolNumber, name).First(&modeOfStudy).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return modeOfStudy, nil
		}

		return modeOfStudy, err
	}

	return modeOfStudy, nil
}

func GetModeOfStudyByNumber(modeOfStudyNumber string, trn ...*gorm.DB) (models.ModeOfStudy, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var modeOfStudy models.ModeOfStudy

	if err := db.Where("number = ?", modeOfStudyNumber).First(&modeOfStudy).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return modeOfStudy, nil
		}

		return modeOfStudy, err
	}

	return modeOfStudy, nil
}
