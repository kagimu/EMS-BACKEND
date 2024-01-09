package models

import "gorm.io/gorm"

type PreviousSchoolClass struct {
	gorm.Model
	PreviousSchoolNumber string `json:"previous_school_number" gorm:"size:36;not null"`
	Year                 string `json:"year" gorm:"size:10;not null"`
	Class                string `json:"class" gorm:"size:10;not null"`
}
