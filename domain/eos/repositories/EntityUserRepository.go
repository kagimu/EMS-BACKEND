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

	if err := db.Preload("Entity").Where("entity_number = ? AND user_number = ? ", entityNumber, userNumber).Find(&entityUser).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entityUser, nil
		}

		return entityUser, err
	}

	return entityUser, nil
}

func FindAllEntityUsers(entityNumber string) ([]models.EntityUser, error) {
	db := config.DB

	var entities []models.EntityUser

	if err := db.Where("entity_number = ? ", entityNumber).Find(&entities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities, nil
		}
		return entities, err
	}
	return entities, nil
}

func DeleteEntityUser(entityUser models.EntityUser, tx ...*gorm.DB) error {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Delete(&entityUser).Error; err != nil {
		return err
	}

	return nil
}
