package models

import (
	"gorm.io/gorm"
)

type StaffBioDataFamily struct {
	gorm.Model
	Number                string `json:"number" gorm:"size:36;unique;index;not null"`
	JobOfferNumber        string `json:"job_offer_number" gorm:"size:36;not null"`
	JobOffer              JobOffer `gorm:"foreignKey:JobOfferNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	
	FatherSurname         string `json:"fatherSurname" gorm:"size:255"`
	FatherOtherNames      string `json:"fatherOtherNames" gorm:"size:255"`
	FatherEmployment      string `json:"fatherEmployment" gorm:"size:255"`
	FatherPhone1          string `json:"fatherPhone1" gorm:"size:255"`
	FatherPhone2          string `json:"fatherPhone2" gorm:"size:255"`
	FatherEmail1          string `json:"fatherEmail1" gorm:"size:255"`
	FatherEmail2          string `json:"fatherEmail2" gorm:"size:255"`
	FatherPhysicalAddress string `json:"fatherPhysicalAddress" gorm:"size:255"`
	FatherLiving          string `json:"fatherLiving" gorm:"size:50"`

	MotherSurname         string `json:"motherSurname" gorm:"size:255"`
	MotherOtherNames      string `json:"motherOtherNames" gorm:"size:255"`
	MotherEmployment      string `json:"motherEmployment" gorm:"size:255"`
	MotherPhone1          string `json:"motherPhone1" gorm:"size:255"`
	MotherPhone2          string `json:"motherPhone2" gorm:"size:255"`
	MotherEmail1          string `json:"motherEmail1" gorm:"size:255"`
	MotherEmail2          string `json:"motherEmail2" gorm:"size:255"`
	MotherPhysicalAddress string `json:"motherPhysicalAddress" gorm:"size:255"`
	MotherLiving          string `json:"motherLiving" gorm:"size:50"`

	BrotherSurname        string `json:"brotherSurname" gorm:"size:255"`
	BrotherOtherNames     string `json:"brotherOtherNames" gorm:"size:255"`
	SisterSurname         string `json:"sisterSurname" gorm:"size:255"`
	SisterOtherNames      string `json:"sisterOtherNames" gorm:"size:255"`

	SpouseSurname         string `json:"spouseSurname" gorm:"size:255"`
	SpouseMaidenName      string `json:"spouseMaidenName" gorm:"size:255"`
	SpouseOtherNames      string `json:"spouseOtherNames" gorm:"size:255"`
	SpouseEmployment      string `json:"spouseEmployment" gorm:"size:255"`
	SpousePhone1          string `json:"spousePhone1" gorm:"size:255"`
	SpousePhone2          string `json:"spousePhone2" gorm:"size:255"`
	SpouseEmail1          string `json:"spouseEmail1" gorm:"size:255"`
	SpouseEmail2          string `json:"spouseEmail2" gorm:"size:255"`
	SpousePhysicalAddress string `json:"spousePhysicalAddress" gorm:"size:255"`
	SpouseLiving          string `json:"spouseLiving" gorm:"size:50"`

	Children []StaffBioDataFamilyChild `json:"children" gorm:"foreignKey:FamilyNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	
	UserNumber string `json:"user_number" gorm:"size:36;not null"`
}
