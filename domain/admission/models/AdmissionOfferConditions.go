package models

import "gorm.io/gorm"

type AdmissionOfferConditions struct {
	gorm.Model
	Number string `json:"number" gorm:"size:36;unique;not null"`
}
