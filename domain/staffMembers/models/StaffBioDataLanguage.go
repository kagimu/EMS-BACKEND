package models

import (
	"gorm.io/gorm"
)

type StaffBioDataLanguage struct {
	gorm.Model
	Number         string                     `json:"number" gorm:"size:36;unique;index;not null"`
	JobOfferNumber string                     `json:"job_offer_number" gorm:"size:36;not null"`
	JobOffer       JobOffer                   `gorm:"foreignKey:JobOfferNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Items          []StaffBioDataLanguageItem `json:"languages" gorm:"foreignKey:LanguageRecordNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserNumber     string                     `json:"user_number" gorm:"size:36;not null"`
}

type StaffBioDataLanguageItem struct {
	gorm.Model
	Number               string `json:"number" gorm:"size:36;unique;index;not null"`
	LanguageRecordNumber string `json:"language_record_number" gorm:"size:36;not null"`
	Language             string `json:"language" gorm:"size:255"`
	Speaking             string `json:"speaking" gorm:"size:50"`
	Reading              string `json:"reading" gorm:"size:50"`
	Writing              string `json:"writing" gorm:"size:50"`
}
