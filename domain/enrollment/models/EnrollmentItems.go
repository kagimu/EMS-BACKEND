package models

import (
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"gorm.io/gorm"
)

type EnrollmentItem struct {
	gorm.Model
	EnrollmentNumber string          `json:"enrollment_number" gorm:"size:36;not null"`
	Enrollment       Enrollment      `gorm:"foreignKey:EnrollmentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ItemNumber       string          `json:"item_number" gorm:"size:36; not null"`
	Item             itemModels.Item `gorm:"foreignKey:ItemNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	RequiredQuantity float64         `json:"required_quantity" gorm:"type:decimal(10,2);default:0"`
	Quantity         float64         `json:"quantity" gorm:"type:decimal(10,2);default:0"`
	Brought          bool            `json:"brought"`
	UserNumber       string          `json:"user_number" gorm:"size:36; not null"`
	User             userModel.User  `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
