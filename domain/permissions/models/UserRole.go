package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

// UserRole represents the relationship between users and roles
type UserRole struct {
	gorm.Model
	UserNumber  string         `json:"user_number" gorm:"size:36;not null"`
	RoleNumber  string         `json:"role_number" gorm:"size:36;not null"`
	AddedBy     string         `json:"added_by" gorm:"size:36;not null"`
	User        authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"` // A user can have many roles
	Role        Role           `gorm:"foreignKey:RoleNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"` // A role can have many users
	AddedByUser authModel.User `gorm:"foreignKey:AddedBy;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`    // A user can add many roles
}
