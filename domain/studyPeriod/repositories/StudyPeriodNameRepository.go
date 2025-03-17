package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

func SaveStudyPeriodName(studyPeriodName models.StudyPeriodName, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&studyPeriodName).Error
}

func FindStudyPeriodName(schoolNumber string, trn ...*gorm.DB) ([]models.StudyPeriodName, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var studyPeriodNames []models.StudyPeriodName

	if err := db.
		Order("name desc").
		Where("school_number = ?", schoolNumber).Find(&studyPeriodNames).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return studyPeriodNames, nil
		}

		return studyPeriodNames, err
	}

	return studyPeriodNames, nil
}

func GetStudyPeriodNameByName(schoolNumber, name string, trn ...*gorm.DB) (models.StudyPeriodName, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var studyPeriodName models.StudyPeriodName

	if err := db.Where("school_number = ? AND name = ?", schoolNumber, name).First(&studyPeriodName).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return studyPeriodName, nil
		}

		return studyPeriodName, err
	}

	return studyPeriodName, nil
}

func GetStudyPeriodNameByNumber(studyPeriodNameNumber string, trn ...*gorm.DB) (models.StudyPeriodName, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var studyPeriodName models.StudyPeriodName

	if err := db.Where("number = ?", studyPeriodNameNumber).First(&studyPeriodName).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.StudyPeriodName{}, nil
		}

		return models.StudyPeriodName{}, err
	}

	return studyPeriodName, nil
}

func RemoveStudyPeriodName(studyPeriodName models.StudyPeriodName, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Delete(&studyPeriodName).Error
}

func GetStudyPeriodNameByStudyPeriodNumber(studyPeriodNumber string, trn ...*gorm.DB) (models.StudyPeriodName, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var studyPeriodName models.StudyPeriodName

	if err := db.
		Joins("JOIN study_periods ON study_periods.study_period_name_number = study_period_names.number").
		Where("study_periods.number = ?", studyPeriodNumber).First(&studyPeriodName).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.StudyPeriodName{}, nil
		}

		return models.StudyPeriodName{}, err
	}

	return studyPeriodName, nil
}
