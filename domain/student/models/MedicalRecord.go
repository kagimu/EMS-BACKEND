package models

import "gorm.io/gorm"

type MedicalRecord struct {
	gorm.Model
	StudentNumber     string `json:"student_number" gorm:"size:36;not null"`
	BloodTypeGroup    string `json:"blood_type_group" gorm:"size:50;not null"`
	Allergies         string `json:"allergies" gorm:"size:50;not null"`
	ExistingCondition string `json:"existing_condition" gorm:"size:100;not null"`
	RegularMedication string `json:"regular_medication" gorm:"size:100;not null"`
	DoctorName        string `json:"doctor_name" gorm:"size:100;not null"`
	DoctorPhone       string `json:"doctor_phone" gorm:"size:100;not null"`
}
