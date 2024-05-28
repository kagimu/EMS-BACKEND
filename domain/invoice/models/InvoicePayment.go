package models

import (
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
	"time"
)

type InvoicePayment struct {
	gorm.Model
	Number        string         `json:"number" gorm:"size:36;unique;not null"`
	Date          time.Time      `json:"date" gorm:"not null"`
	Amount        float64        `json:"amount"`
	InvoiceNumber string         `json:"invoice_number" gorm:"size:36;not null"`
	Invoice       Invoice        `gorm:"foreignKey:InvoiceNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber    string         `json:"user_number" gorm:"size:36;not null;comment:'Received or recorded by'"`
	User          userModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
