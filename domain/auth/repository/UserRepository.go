package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

func GetUserByNumber(userNumber string) (models.User, error) {

	db := config.DB

	var user models.User

	if err := db.Where("number = ?", userNumber).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, nil
		}
		return user, err
	}

	return user, nil
}
