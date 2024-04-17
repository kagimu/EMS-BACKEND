package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

// RolePermission stores the relationship between roles and permissions
type RolePermission struct {
	gorm.Model
	ID            uint
	RoleNumber    string          `gorm:"size:36;not null" json:"role_number"`
	Permission    string          `gorm:"size:36;not null" json:"permission"`
	AddedByNumber string          `gorm:"size:36;not null" json:"added_by_number"`
	AddedBy       authModels.User `gorm:"foreignKey:AddedByNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Role          Role            `gorm:"foreignKey:RoleNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
