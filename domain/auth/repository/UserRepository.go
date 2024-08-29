package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	permissionModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
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

func GetUserByNumber(userNumber string, txn ...*gorm.DB) (models.User, error) {

	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var user models.User

	if err := db.Where("number = ?", userNumber).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, nil
		}
		return user, err
	}

	return user, nil
}

func FetchAllPermissionsByUser(userNumber string) ([]string, error) {
	db := config.DB
	//prepare roles and permissions
	//get user roles
	var userRoles []permissionModels.UserRole
	permissions := make([]string, 0)

	if err := db.Where("user_number = ?", userNumber).Find(&userRoles).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return permissions, nil
		}
		return permissions, err

	}

	for _, item := range userRoles {
		//get role permissions
		var rolePermissions []permissionModels.RolePermission
		if err := db.Where("role_number = ?", item.RoleNumber).Find(&rolePermissions).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return permissions, nil
			}
			return permissions, err
		}

		for _, permission := range rolePermissions {
			permissions = append(permissions, permission.Permission)
		}

	}
	return permissions, nil

}
