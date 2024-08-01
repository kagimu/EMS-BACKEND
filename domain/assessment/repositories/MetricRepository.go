package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"gorm.io/gorm"
)

func GetMetricByNumber(metricNumber string, tx ...*gorm.DB) (models.Metric, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var metric models.Metric

	if err := db.Where("number = ?", metricNumber).First(&metric).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Metric{}, nil
		}

		return models.Metric{}, err

	}

	return metric, nil
}
