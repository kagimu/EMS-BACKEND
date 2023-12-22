package models

import "gorm.io/gorm"

type PreviousSchoolClass struct {
	gorm.Model
	Number      string `json:"number" gorm:"size:36;unique;not null"`
	ClassNumber string `json:"class_number" gorm:"size:36;not null"`
}
