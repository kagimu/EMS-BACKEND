package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"gorm.io/gorm"
)

func FindEntityByNumber(entityNumber string) (models.Entity, error) {
	db := config.DB
	var entity models.Entity

	if err := db.Where("number = ?", entityNumber).First(&entity).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, nil
		}

		return entity, err
	}

	return entity, nil
}
