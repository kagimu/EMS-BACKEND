package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/models"
	"gorm.io/gorm"
)

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

func FindPromotionSettingValueBySettingNumber(promotionSettingNumber string) (models.PromotionSettingValue, error) {
	var promotionSettings models.PromotionSettingValue
	db := config.DB

	if err := db.Where("promotion_setting_number = ?", promotionSettingNumber).Find(&promotionSettings).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.PromotionSettingValue{}, nil
		}

		return models.PromotionSettingValue{}, err
	}

	return promotionSettings, nil
}
