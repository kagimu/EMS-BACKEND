package models

import (
	"gorm.io/gorm"
)

type StaffBioDataEducation struct {
	gorm.Model
	Number         string   `json:"number" gorm:"size:36;unique;index;not null"`
	JobOfferNumber string   `json:"job_offer_number" gorm:"size:36;not null"`
	JobOffer       JobOffer `gorm:"foreignKey:JobOfferNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`

	Items      []StaffBioDataEducationItem `json:"items" gorm:"foreignKey:EducationNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserNumber string                      `json:"user_number" gorm:"size:36;not null"`
}
