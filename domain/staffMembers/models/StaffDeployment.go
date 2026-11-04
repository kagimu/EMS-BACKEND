package models

import (
	"gorm.io/gorm"
)

type StaffDeployment struct {
	gorm.Model
	Number         string   `json:"number" gorm:"size:36;unique;index;not null"`
	JobOfferNumber string   `json:"job_offer_number" gorm:"size:36;not null"`
	JobOffer       JobOffer `gorm:"foreignKey:JobOfferNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StaffIdNumber  string   `json:"staff_id_number" gorm:"size:100;unique;not null"`
	Status         string   `json:"status" gorm:"size:50;default:'Deployed'"`
	UserNumber     string   `json:"user_number" gorm:"size:36;not null"`
}
