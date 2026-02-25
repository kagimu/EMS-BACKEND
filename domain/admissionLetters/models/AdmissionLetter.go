package models

import (
	"gorm.io/gorm"
)

type AdmissionLetter struct {
	gorm.Model
	Number              string `json:"number" gorm:"size:36;unique;not null"`
	SchoolNumber        string `json:"school_number" gorm:"size:36;not null"`
	StudentNumber       string `json:"student_number" gorm:"size:36;not null"`
	CampusNumber        string `json:"campus_number" gorm:"size:36;not null"`
	CampusClassNumber   string `json:"campus_class_number" gorm:"size:36;not null"`
	Phrase              string `json:"phrase" gorm:"type:text"`
	PrincipalSignature  string `json:"principal_signature" gorm:"size:255"`
}

func (AdmissionLetter) TableName() string {
	return "admission_letters"
}
