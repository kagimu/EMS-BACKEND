package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetUnitOfMeasureByNumber(unitOfMeasureNumber string) (models.UnitOfMeasure, error) {

	db := config.DB

	var unitOfMeasure models.UnitOfMeasure

	if err := db.Where("number = ? ", unitOfMeasureNumber).First(&unitOfMeasure).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return unitOfMeasure, nil
		}

		return unitOfMeasure, err

	}

	return unitOfMeasure, nil
}
