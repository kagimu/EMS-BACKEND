package implementation

import (
	"errors"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	"github.com/jinzhu/gorm"
)

func GetUser(number string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("number = ?", number).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
