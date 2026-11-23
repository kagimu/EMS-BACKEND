package models

import (
	"gorm.io/gorm"
)

type StaffBioDataFamilyChild struct {
	gorm.Model
	Number       string `json:"number" gorm:"size:36;unique;index;not null"`
	FamilyNumber string `json:"family_number" gorm:"size:36;not null"`
	Surname      string `json:"surname" gorm:"size:255"`
	OtherNames   string `json:"otherNames" gorm:"size:255"`
	Dob          string `json:"dob" gorm:"size:50"`
}
