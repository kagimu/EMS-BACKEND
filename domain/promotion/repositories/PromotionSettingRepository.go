package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/models"
	"gorm.io/gorm"
)

func FindPromotionSettingsByStudyPeriod(studyPeriodNumber, currentCampusClassNumber, checkType string) (models.PromotionSetting, error) {
	var promotionSettings models.PromotionSetting
	db := config.DB

	if err := db.Where("study_period_number = ? AND previous_campus_class_number = ? AND type = ? AND status = ?", studyPeriodNumber, currentCampusClassNumber, checkType, constants.StatusActive).Find(&promotionSettings).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.PromotionSetting{}, nil
		}

		return models.PromotionSetting{}, err
	}

	return promotionSettings, nil
}

func FindPromotionSettingValueBySettingNumber(promotionSettingNumber string, trn ...*gorm.DB) (models.PromotionSettingValue, error) {
	var promotionSettings models.PromotionSettingValue
	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Where("promotion_setting_number = ?", promotionSettingNumber).Find(&promotionSettings).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.PromotionSettingValue{}, nil
		}

		return models.PromotionSettingValue{}, err
	}

	return promotionSettings, nil
}

func FindAllPromotionSettings(studyPeriodNumber, nextCampusClassNumber, settingsType string) ([]models.PromotionSetting, error) {
	var promotionSettings []models.PromotionSetting
	db := config.DB

	if err := db.Scopes(
		filters.ByAnyField("type", settingsType),
		filters.ByAnyField("next_campus_class_number", nextCampusClassNumber),
		filters.ByAnyField("study_period_number", studyPeriodNumber),
	).Find(&promotionSettings).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.PromotionSetting{}, nil
		}

		return []models.PromotionSetting{}, err
	}

	return promotionSettings, nil
}

func SavePromotionSettings(promotionSettings models.PromotionSetting) error {
	db := config.DB

	if err := db.Save(&promotionSettings).Error; err != nil {
		return err
	}

	return nil
}

func SavePromotionSettingValue(promotionSettingValue models.PromotionSettingValue, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Save(&promotionSettingValue).Error; err != nil {
		return err
	}

	return nil
}

func SavePromotionAssessmentAreaValue(promotionAssessmentAreaValue models.PromotionAssessmentAreaValue, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Save(&promotionAssessmentAreaValue).Error; err != nil {
		return err
	}

	return nil
}

func FindPromotionAssessmentAreaValueBySettingsNumberAndAssessmentAreaNumber(promotionSettingNumber, assessmentAreaNumber string, trn ...*gorm.DB) (models.PromotionAssessmentAreaValue, error) {
	var promotionSettings models.PromotionAssessmentAreaValue
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Where("promotion_setting_number = ? AND assessment_area_number", promotionSettingNumber, assessmentAreaNumber).Find(&promotionSettings).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.PromotionAssessmentAreaValue{}, nil
		}

		return models.PromotionAssessmentAreaValue{}, err
	}

	return promotionSettings, nil
}

func FindPromotionSettingByNumber(promotionSettingNumber string, trn ...*gorm.DB) (models.PromotionSetting, error) {
	var promotionSettings models.PromotionSetting
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Where("number = ?", promotionSettingNumber).Find(&promotionSettings).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.PromotionSetting{}, nil
		}

		return models.PromotionSetting{}, err
	}

	return promotionSettings, nil
}

func SavePromotion(promotion models.Promotion, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Save(&promotion).Error; err != nil {
		return err
	}

	return nil
}

func FindPromotionByReportSummaryNumber(reportSummaryNumber string, trn ...*gorm.DB) (models.Promotion, error) {
	var promotion models.Promotion
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Where("report_summary_number = ?", reportSummaryNumber).Find(&promotion).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Promotion{}, nil
		}

		return models.Promotion{}, err
	}

	return promotion, nil
}
