package models

import (
	"gorm.io/gorm"
)

type StaffBioDataPersonal struct {
	gorm.Model
	Number             string   `json:"number" gorm:"size:36;unique;index;not null"`
	JobOfferNumber     string   `json:"job_offer_number" gorm:"size:36;not null"`
	JobOffer           JobOffer `gorm:"foreignKey:JobOfferNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Surname            string   `json:"surname" gorm:"size:255;not null"`
	SecondName         string   `json:"secondName" gorm:"size:255;not null"`
	OtherNames         string   `json:"otherNames" gorm:"size:255;not null"`
	DateCompleted      string   `json:"dateCompleted" gorm:"size:50;not null"`
	DateOfBirth        string   `json:"dateOfBirth" gorm:"size:50;not null"`
	Sex                string   `json:"sex" gorm:"size:50;not null"`
	BirthCertificateNo string   `json:"birthCertificateNo" gorm:"size:255;not null"`
	Religion           string   `json:"religion" gorm:"size:255;not null"`
	NationalIdNo       string   `json:"nationalIdNo" gorm:"size:255;not null"`
	CountryOfBirth     string   `json:"countryOfBirth" gorm:"size:255;not null"`
	HomeCountry        string   `json:"homeCountry" gorm:"size:255;not null"`
	HomeDistrict       string   `json:"homeDistrict" gorm:"size:255;not null"`
	HomeCounty         string   `json:"homeCounty" gorm:"size:255;not null"`
	HomeSubCounty      string   `json:"homeSubCounty" gorm:"size:255;not null"`
	HomeParish         string   `json:"homeParish" gorm:"size:255;not null"`
	HomeVillage        string   `json:"homeVillage" gorm:"size:255;not null"`
	VillageOfResidence string   `json:"villageOfResidence" gorm:"size:255;not null"`
	TaxIdNo            string   `json:"taxIdNo" gorm:"size:255;not null"`
	UserNumber         string   `json:"user_number" gorm:"size:36;not null"`
	Photo              string   `json:"photo" gorm:"size:255"`
}
