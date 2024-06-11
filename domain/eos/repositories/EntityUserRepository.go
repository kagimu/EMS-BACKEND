package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"gorm.io/gorm"
)

func FindEntityUser(userNumber string, entityNumber string) (models.EntityUser, error) {

	db := config.DB

	var entityUser models.EntityUser

	if err := db.Where("entity_number = ? AND user_number = ? ", entityNumber, userNumber).Find(&entityUser).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entityUser, nil
		}

		return entityUser, err
	}

	return entityUser, nil
}
