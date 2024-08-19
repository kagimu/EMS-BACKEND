package models

import "gorm.io/gorm"

type CampusItem struct {
	gorm.Model
	Number       string `json:"number" gorm:"size:36;not null"`
	CampusNumber string `json:"campus_number" gorm:"size:36;"`
	ItemNumber   string `json:"item_number" gorm:"size:36;"`
}
