package models

import (
	"gorm.io/gorm"
)

type StaffBioDataContact struct {
	gorm.Model
	Number            string   `json:"number" gorm:"size:36;unique;index;not null"`
	JobOfferNumber    string   `json:"job_offer_number" gorm:"size:36;not null"`
	JobOffer          JobOffer `gorm:"foreignKey:JobOfferNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	MainEmail         string   `json:"mainEmail" gorm:"size:255;not null"`
	SecondEmail       string   `json:"secondEmail" gorm:"size:255"`
	AltEmail          string   `json:"altEmail" gorm:"size:255"`
	MainPhone         string   `json:"mainPhone" gorm:"size:255;not null"`
	SecondPhone       string   `json:"secondPhone" gorm:"size:255"`
	AltPhone          string   `json:"altPhone" gorm:"size:255"`
	PostMailAddress   string   `json:"postMailAddress" gorm:"size:255"`
	PhysicalAddress   string   `json:"physicalAddress" gorm:"size:255;not null"`
	UserNumber        string   `json:"user_number" gorm:"size:36;not null"`
}
