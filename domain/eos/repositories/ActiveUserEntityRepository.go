package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"gorm.io/gorm"
)

func FindActiveUserEntityByUserNumber(userNumber string) (models.ActiveUserEntity, error) {

	db := config.DB

	var activeUserEntity models.ActiveUserEntity

	if err := db.Preload("Entity").Where("user_number = ?", userNumber).First(&activeUserEntity).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return activeUserEntity, nil
		}

		return activeUserEntity, err
	}

	return activeUserEntity, nil

}

func SaveActiveUserEntity(activeUserEntity models.ActiveUserEntity) error {

	db := config.DB

	if err := db.Save(&activeUserEntity).Error; err != nil {
		return err
	}

	return nil
}
