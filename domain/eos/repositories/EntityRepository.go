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

func FindEntityByName(entityName string) (models.Entity, error) {

	db := config.DB
	var entity models.Entity

	if err := db.Where("name = ?", entityName).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, nil
		}
		return entity, err
	}

	return entity, nil
}

func SaveEntity(entity models.Entity) error {

	db := config.DB

	if err := db.Save(&entity).Error; err != nil {
		return err
	}

	return nil
}

func DeleteEntity(entity models.Entity, tx ...*gorm.DB) error {
	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Delete(&entity).Error; err != nil {
		return err
	}

	return nil
}
