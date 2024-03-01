package models

import "gorm.io/gorm"

type ExaminationType struct {
	gorm.Model
	Number      string `json:"number" gorm:"size:36;unique;not null"`
	Name        string `json:"name" gorm:"size:20;"`
	Description string `json:"description" gorm:"type:text;"`
}
