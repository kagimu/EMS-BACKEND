package models

import (
	"gorm.io/gorm"
)

type InvoicePayment struct {
	gorm.Model
	Number        string `json:"number" gorm:"size:36;unique;not null"`
	Date          string `json:"date"`
	InvoiceNumber string `json:"invoice_number" gorm:"size:36;not null"`
	Amount        int    `json:"amount"`
	Balance       int    `json:"-"`
	Payer         string `json:"payer" gorm:"size:50"`
}
