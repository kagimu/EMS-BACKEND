package models

import (
	"gorm.io/gorm"
)

type StaffBioDataMedical struct {
	gorm.Model
	Number              string   `json:"number" gorm:"size:36;unique;index;not null"`
	JobOfferNumber      string   `json:"job_offer_number" gorm:"size:36;not null"`
	JobOffer            JobOffer `gorm:"foreignKey:JobOfferNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Height              string   `json:"height" gorm:"size:255"`
	Weight              string   `json:"weight" gorm:"size:255"`
	Complexion          string   `json:"complexion" gorm:"size:255"`
	BloodType           string   `json:"bloodType" gorm:"size:50"`
	Allergies           string   `json:"allergies" gorm:"type:text"`
	MedicalConditions   string   `json:"medicalConditions" gorm:"type:text"`
	RegularMedications  string   `json:"regularMedications" gorm:"type:text"`
	PhysicalImpairments string   `json:"physicalImpairments" gorm:"type:text"`
	DoctorName          string   `json:"doctorName" gorm:"size:255"`
	DoctorPhone         string   `json:"doctorPhone" gorm:"size:255"`
	DoctorFacility      string   `json:"doctorFacility" gorm:"size:255"`
	UserNumber          string   `json:"user_number" gorm:"size:36;not null"`
}
