package models

import "gorm.io/gorm"

type CampusServices struct {
	gorm.Model
	Number        string `json:"number" gorm:"size:36;not null"`
	CampusNumber  string `json:"campus_number" gorm:"size:36;"`
	ServiceNumber string `json:"service_number" gorm:"size:36;"`
}
