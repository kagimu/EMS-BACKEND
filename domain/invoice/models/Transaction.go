package models

import (
	"time"
)

type Transaction struct {
	Number          string    `json:"number" gorm:"size:36;primaryKey;not null"`
	Date            time.Time `json:"date" gorm:"not null"`
	AmountInCredit  float64   `json:"amount_in_credit"`
	Status          string    `json:"status" gorm:"size:36;not null;default:ACTIVE"`
	ClientReference string    `json:"client_reference" gorm:"size:62"`
	Description     string    `json:"description" gorm:"size:500"`
	StatementNumber string    `json:"statement_number" gorm:"size:36;not null;unique"`
	Statement       Statement `gorm:"foreignKey:StatementNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
