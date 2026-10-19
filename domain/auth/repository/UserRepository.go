package repository

import (
	"errors"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
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

// FetchAllPermissionsByUser returns every effective permission the user holds
// across all organisations. It is the installation-level view used only by the
// explicit control-plane authorities (licence registry and audit registry); it
// still applies the role-status and assignment-expiry rules.
func FetchAllPermissionsByUser(userNumber string) ([]string, error) {
	return fetchPermissions(userNumber, "")
}

// FetchPermissionsByUserAndEntity returns the effective permissions the user
// holds through ACTIVE, unexpired roles that belong to the given organisation.
//
// This is the organisation-scoped view used by business-request authorisation
// and by role administration to decide what an actor may delegate.
func FetchPermissionsByUserAndEntity(userNumber string, entityNumber string) ([]string, error) {
	return fetchPermissions(userNumber, entityNumber)
}

// fetchPermissions applies the shared status and expiry rules. An empty
// entityNumber keeps the installation-wide view; a non-empty value adds the
// organisation predicate.
func fetchPermissions(userNumber string, entityNumber string) ([]string, error) {
	permissions := make([]string, 0)
	if userNumber == "" {
		return permissions, nil
	}

	query := config.DB.
		Table("user_roles").
		Select("DISTINCT role_permissions.permission").
		Joins("JOIN roles ON roles.number = user_roles.role_number AND roles.deleted_at IS NULL").
		Joins("JOIN role_permissions ON role_permissions.role_number = roles.number AND role_permissions.deleted_at IS NULL").
		Where("user_roles.user_number = ?", userNumber).
		Where("user_roles.deleted_at IS NULL").
		Where("user_roles.expiry_date >= ?", time.Now()).
		Where("roles.status = ?", constants.StatusActive)

	if entityNumber != "" {
		query = query.Where("roles.entity_number = ?", entityNumber)
	}

	err := query.Scan(&permissions).Error
	return permissions, err
}
