package implementation

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
)

func GetUser(number string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("number = ?", number).Take(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
