package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

func SaveClassStudyPeriod(classStudyPeriod models.ClassStudyPeriod, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&classStudyPeriod).Error
}

func FindClassStudyPeriods(classNumber, studyPeriodNumber string, trn ...*gorm.DB) ([]models.ClassStudyPeriod, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var studyPeriods []models.ClassStudyPeriod

	if err := db.
		Scopes(filters.ByAnyField("class_number", classNumber), filters.ByAnyField("study_period_number", studyPeriodNumber)).
		Find(&studyPeriods).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return studyPeriods, nil
		}

		return studyPeriods, err
	}

	return studyPeriods, nil
}

func DeleteClassStudyPeriod(classStudyPeriod models.ClassStudyPeriod, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Delete(&classStudyPeriod).Error
}
