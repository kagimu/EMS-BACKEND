package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/activityreporting/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func FindIndicatorByName(indicatorName string, entityNumber string) (models.Indicator, error) {
	db := config.DB
	var indicator models.Indicator

	if err := db.Where("name = ? AND entity_number = ?", indicatorName, entityNumber).First(&indicator).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return indicator, nil
		}

		return indicator, err
	}

	return indicator, nil
}

func FindIndicatorByCode(indicatorCode int, entityNumber string) (models.Indicator, error) {
	db := config.DB
	var indicator models.Indicator

	if err := db.Where("code = ? AND entity_number = ?", indicatorCode, entityNumber).First(&indicator).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return indicator, nil
		}
		return indicator, err
	}
	return indicator, nil
}

func SaveIndicator(indicator models.Indicator) error {
	if err := config.DB.Save(&indicator).Error; err != nil {
		return err
	}
	return nil
}

func GetIndicators(entityNumber string, indicatorName string, status string) ([]models.Indicator, error) {

	db := config.DB

	var indicators = make([]models.Indicator, 0)

	if err := db.Scopes(filters.ByStatus(status), filters.ByName(indicatorName)).
		Where("entity_number = ?", entityNumber).
		Find(&indicators).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return indicators, nil
		}

		return indicators, err
	}

	return indicators, nil
}
