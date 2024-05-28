package repository

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
)

func FindRoleByNumber(roleNumber string) (*models.Role, error) {
	var role models.Role
	if err := config.DB.Where("number = ?", roleNumber).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
