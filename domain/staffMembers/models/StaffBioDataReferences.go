package models

import (
	"gorm.io/gorm"
)

type StaffBioDataReferences struct {
	gorm.Model
	Number         string   `json:"number" gorm:"size:36;unique;index;not null"`
	JobOfferNumber string   `json:"job_offer_number" gorm:"size:36;not null"`
	JobOffer       JobOffer `gorm:"foreignKey:JobOfferNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`

	Guarantor1Name         string `json:"guarantor1_name" gorm:"size:255"`
	Guarantor1Address      string `json:"guarantor1_address" gorm:"size:255"`
	Guarantor1Email        string `json:"guarantor1_email" gorm:"size:255"`
	Guarantor1Phone        string `json:"guarantor1_phone" gorm:"size:255"`
	Guarantor1Relationship string `json:"guarantor1_relationship" gorm:"size:255"`

	Guarantor2Name         string `json:"guarantor2_name" gorm:"size:255"`
	Guarantor2Address      string `json:"guarantor2_address" gorm:"size:255"`
	Guarantor2Email        string `json:"guarantor2_email" gorm:"size:255"`
	Guarantor2Phone        string `json:"guarantor2_phone" gorm:"size:255"`
	Guarantor2Relationship string `json:"guarantor2_relationship" gorm:"size:255"`

	NextOfKinName    string `json:"nextOfKin_name" gorm:"size:255"`
	NextOfKinAddress string `json:"nextOfKin_address" gorm:"size:255"`
	NextOfKinEmail   string `json:"nextOfKin_email" gorm:"size:255"`
	NextOfKinPhone   string `json:"nextOfKin_phone" gorm:"size:255"`

	Referee1Name         string `json:"referee1_name" gorm:"size:255"`
	Referee1Address      string `json:"referee1_address" gorm:"size:255"`
	Referee1Email        string `json:"referee1_email" gorm:"size:255"`
	Referee1Phone        string `json:"referee1_phone" gorm:"size:255"`
	Referee1Relationship string `json:"referee1_relationship" gorm:"size:255"`
	Referee1YearKnown    string `json:"referee1_yearKnown" gorm:"size:50"`

	Referee2Name         string `json:"referee2_name" gorm:"size:255"`
	Referee2Address      string `json:"referee2_address" gorm:"size:255"`
	Referee2Email        string `json:"referee2_email" gorm:"size:255"`
	Referee2Phone        string `json:"referee2_phone" gorm:"size:255"`
	Referee2Relationship string `json:"referee2_relationship" gorm:"size:255"`
	Referee2YearKnown    string `json:"referee2_yearKnown" gorm:"size:50"`

	UserNumber string `json:"user_number" gorm:"size:36;not null"`
}
