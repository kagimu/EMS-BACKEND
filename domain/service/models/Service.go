package models

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	Number        string `json:"number" gorm:"size:36;unique;not null"`
	Description   string `json:"description" gorm:"type:text"`
	Name          string `json:"name" gorm:"size:255"`
	StudentNumber string `json:"student_number" gorm:"size:36;index;not null"`
	Fee           uint   `json:"fee"`
	Status        string `json:"status" gorm:"size:36"`
}
