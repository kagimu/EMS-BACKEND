package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepo interface {
	CreateUser(ctx *gin.Context, user models.User) error
	GetUserByNumber(ctx *gin.Context, userNumber string) (models.User, error)
}
type SqlUserRepo struct {
	db *gorm.DB
}

func NewSqlUserRepository(db *gorm.DB) *SqlUserRepo {
	return &SqlUserRepo{db: db}
}

func (s *SqlUserRepo) CreateUser(user models.User) error {
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil

}

func (s *SqlUserRepo) GetUserByNumber(userNumber string) (models.User, error) {

	var user models.User

	if err := s.db.Where("number = ?", userNumber).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, nil
		}
		return user, err
	}

	return user, nil

}

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
