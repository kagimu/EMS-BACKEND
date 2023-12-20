package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Application struct {
	gorm.Model
	Number                string      `json:"number" gorm:"size:36;unique;not null"`
	Status                string      `gorm:"type:varchar(255);not null;default:'DRAFT'" json:"status"`
	ApplicationFeePayment bool        `json:"application_fee_payment" gorm:"not null;default:false"`
	UserNumber            string      `json:"user_number" gorm:"size:36;index;not null"`
	User                  models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
