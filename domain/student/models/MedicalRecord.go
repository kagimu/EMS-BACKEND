package models

import (
	"time"
)

type MedicalRecord struct {
	StudentNumber     string    `json:"student_number" gorm:"size:36;not null"`
	BloodTypeGroup    string    `json:"blood_type_group" gorm:"size:50;not null"`
	Allergies         string    `json:"allergies" gorm:"size:50;not null"`
	ExistingCondition string    `json:"existing_condition" gorm:"size:100;not null"`
	RegularMedication string    `json:"regular_medication" gorm:"size:100;not null"`
	DoctorName        string    `json:"doctor_name" gorm:"size:100;not null"`
	DoctorPhone       string    `json:"doctor_phone" gorm:"size:100;not null"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
