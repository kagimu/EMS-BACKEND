package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Number       string `json:"number" gorm:"size:36;unique;not null"`
	Name         string `json:"name" gorm:"size:20;not null"`
	Description  string `json:"description" gorm:"size:100"`
	ClassNumber  string `gorm:"size:36;" json:"class_number"`
	CampusNumber string `gorm:"size:36;" json:"campus_number"`
}
