package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	Number            string      `json:"number" gorm:"size:36;unique;not null"`
	Status            string      `json:"status"`
	Amount            int         `json:"amount"`
	Details           string      `json:"details"`
	ApplicationNumber string      `json:"application_number"`
	Method            string      `json:"method"`
	Date              string      `json:"date"`
	SlipSerialNumber  string      `json:"serial_number"`
	ApprovedBy        string      `json:"approved_by"`
	UserNumber        string      `json:"user_number" gorm:"size:36;index;not null"`
	User              models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
