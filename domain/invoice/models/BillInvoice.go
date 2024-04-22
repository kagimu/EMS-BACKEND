package models

import "gorm.io/gorm"

type BillInvoice struct {
	gorm.Model
	Number       string `json:"number" gorm:"size:36;unique;not null"`
	CampusNumber string `json:"campus_number" gorm:"size:36;not null"`
}
