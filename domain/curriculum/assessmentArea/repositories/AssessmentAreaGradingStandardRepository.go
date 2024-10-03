package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	"gorm.io/gorm"
)

func FindAssessmentAreaGradingStandardsByGradingStandardNumber(schoolNumber, assessmentAreaNumber, gradingStandardNumber string, trn ...*gorm.DB) ([]models.AssessmentAreaGradingStandard, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var gradingStandards []models.AssessmentAreaGradingStandard

	if err := db.Scopes(filters.ByGradingStandardNumber(gradingStandardNumber), filters.BySchoolNumber(schoolNumber), filters.ByAssessmentAreaNumber(assessmentAreaNumber)).Find(&gradingStandards).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gradingStandards, nil
		}

		return gradingStandards, err
	}

	return gradingStandards, nil
}

func SaveAssessmentAreaGradingStandard(assessmentAreaGradingStandard models.AssessmentAreaGradingStandard, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&assessmentAreaGradingStandard).Error
}

func GetAssessmentAreaGradingStandardsByCodeAndSchoolNumber(code int64, schoolNumber string, trn ...*gorm.DB) (models.AssessmentAreaGradingStandard, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var gradingStandards models.AssessmentAreaGradingStandard

	if err := db.Where("school_number = ? AND code = ?", schoolNumber, code).First(&gradingStandards).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gradingStandards, nil
		}

		return gradingStandards, err
	}

	return gradingStandards, nil
}
