package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Number             string      `json:"number" gorm:"size:36;unique;not null"`
	Name               string      `json:"name"`
	Description        string      `json:"description" gorm:"type:text"`
	SerialCode         string      `json:"serial_code"`
	Prefix             string      `json:"prefix"`
	DiscountPercentage uint        `json:"discount_percentage"`
	DiscountAmount     uint        `json:"discount_amount"`
	Cost               uint        `json:"cost"`
	UserNumber         string      `json:"user_number" gorm:"size:36;not null"`
	User               models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
