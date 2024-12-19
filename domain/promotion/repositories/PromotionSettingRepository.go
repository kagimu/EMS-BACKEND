package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
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
	).Find(promotionSettings).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.PromotionSetting{}, nil
		}

		return []models.PromotionSetting{}, err
	}

	return promotionSettings, nil
}
