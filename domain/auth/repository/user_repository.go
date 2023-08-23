package repository

import (
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
)

type UserRepository interface {
	SaveUser(*models.User) (*models.User, map[string]string)
	GetUser(number string) (*models.User, error)
	GetUsers() ([]models.User, error)
	GetUserByEmailAndPassword(*models.User) (*models.User, map[string]string)
}
