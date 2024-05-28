package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Number      string           `json:"number" gorm:"size:36;unique;not null"`
	Name        string           `json:"name" gorm:"size:20;not null"`
	Description string           `json:"description" gorm:"size:100"`
	Status      string           `json:"status" gorm:"size:20;"`
	UserNumber  string           `json:"user_number" gorm:"size:36;index;not null"`
	User        authModel.User   `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Permissions []RolePermission `json:"permissions" gorm:"foreignKey:RoleNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
