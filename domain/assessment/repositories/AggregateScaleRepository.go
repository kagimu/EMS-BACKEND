package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func SaveAggregateScale(aggregateScale models.AggregateScale, trn ...*gorm.DB) error {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&aggregateScale).Error
}

func FindAggregateScaleByScaleNumber(scaleNumber string, trn ...*gorm.DB) ([]models.AggregateScale, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var aggregateScales []models.AggregateScale

	if err := db.Where("scale_number = ?", scaleNumber).Find(&aggregateScales).Error; err != nil {

		if errors.Is(gorm.ErrRecordNotFound, err) {
			return aggregateScales, nil
		}

		return aggregateScales, err
	}

	return aggregateScales, nil
}

func GetAggregateScaleByScaleNumberAndAssessmentGradeNumber(scaleNumber, assessmentGradeNumber string, trn ...*gorm.DB) (models.AggregateScale, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var aggregateScales models.AggregateScale

	if err := db.Where("scale_number = ? AND assessment_grade_number = ?", scaleNumber, assessmentGradeNumber).Find(&aggregateScales).Error; err != nil {

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

func GetAggregateScaleByScaleNumberAndWeight(scaleNumber string, weight float64, trn ...*gorm.DB) (models.AggregateScale, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var aggregateScales models.AggregateScale

	if err := db.Where("scale_number = ? AND weight = ?", scaleNumber, weight).Find(&aggregateScales).Error; err != nil {

		if errors.Is(gorm.ErrRecordNotFound, err) {
			return aggregateScales, nil
		}

		return aggregateScales, err
	}

	return aggregateScales, nil
}

func GetAggregateScaleByScaleNumberAndRemarks(scaleNumber, remarks string, trn ...*gorm.DB) (models.AggregateScale, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var aggregateScales models.AggregateScale

	if err := db.Where("scale_number = ? AND remarks = ?", scaleNumber, remarks).Find(&aggregateScales).Error; err != nil {

		if errors.Is(gorm.ErrRecordNotFound, err) {
			return aggregateScales, nil
		}

		return aggregateScales, err
	}

	return aggregateScales, nil
}

func GetAggregateScaleByScaleNumberAndCode(scaleNumber string, code int64, trn ...*gorm.DB) (models.AggregateScale, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var aggregateScales models.AggregateScale

	if err := db.Where("scale_number = ? AND code = ?", scaleNumber, code).Find(&aggregateScales).Error; err != nil {

		if errors.Is(gorm.ErrRecordNotFound, err) {
			return aggregateScales, nil
		}

		return aggregateScales, err
	}

	return aggregateScales, nil
}

func FindAggregateScaleByClassNumberAndGradingStandardNumber(classNumber, gradingStandardNumber string, trn ...*gorm.DB) ([]models.AggregateScale, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var aggregateScales []models.AggregateScale

	if err := db.
		Joins("join scales on scales.number = aggregate_scales.scale_number").
		Joins("join grading_standard_class_scale on grading_standard_class_scale.scale_number = scales.number").
		Scopes(filters.ByAnyField("grading_standard_class_scale.class_number", classNumber), filters.ByAnyField("grading_standard_class_scale.grading_standard_number", gradingStandardNumber)).
		Find(&aggregateScales).Error; err != nil {

		if errors.Is(gorm.ErrRecordNotFound, err) {
			return aggregateScales, nil
		}

		return aggregateScales, err
	}

	return aggregateScales, nil
}
