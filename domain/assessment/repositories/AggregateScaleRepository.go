package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	"gorm.io/gorm"
)

func SaveAggregateScale(aggregateScale models.AggregateScale, trn ...*gorm.DB) error {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&aggregateScale).Error
}

func FindAggregateScaleByGradingStandardNumber(gradingStandardNumber string, trn ...*gorm.DB) ([]models.AggregateScale, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var aggregateScales []models.AggregateScale

	if err := db.Where("grading_standard_number = ?", gradingStandardNumber).Find(&aggregateScales).Error; err != nil {

		if errors.Is(gorm.ErrRecordNotFound, err) {
			return aggregateScales, nil
		}

		return aggregateScales, err
	}

	return aggregateScales, nil
}

func GetAggregateScaleByGradingStandardNumberAndGradeNumber(gradingStandardNumber string, gradeNumber string, trn ...*gorm.DB) (models.AggregateScale, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var aggregateScales models.AggregateScale

	if err := db.Where("grading_standard_number = ? AND grade_number = ?", gradingStandardNumber, gradeNumber).Find(&aggregateScales).Error; err != nil {

		if errors.Is(gorm.ErrRecordNotFound, err) {
			return aggregateScales, nil
		}

		return aggregateScales, err
	}

	return aggregateScales, nil
}

func RemoveAggregateScale(aggregateScale models.AggregateScale, trn ...*gorm.DB) error {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Delete(&aggregateScale).Error
}
