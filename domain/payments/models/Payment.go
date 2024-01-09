package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	Number            string `json:"number" gorm:"size:36;unique;not null"`
	Status            string `json:"status"`
	Amount            int    `json:"amount"`
	Details           string `json:"details"`
	ApplicationNumber string `json:"application_number"`
	Method            string `json:"method"`
	Date              string `json:"date"`
}
