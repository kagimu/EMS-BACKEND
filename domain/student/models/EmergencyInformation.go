package models

import "time"

type EmergencyInformation struct {
	StudentNumber                 string    `json:"student_number" gorm:"size:36;not null"`
	ContactPerson                 string    `json:"contact_person" gorm:"size:50;not null"`
	ContactPersonPhone            string    `json:"contact_person_phone" gorm:"size:50;not null"`
	ContactPersonAlternative      string    `json:"contact_person_alternative" gorm:"size:50;not null"`
	ContactPersonAlternativePhone string    `json:"contact_person_alternative_phone" gorm:"size:50;not null"`
	ContactDoctorName             string    `json:"contact_doctor_name" gorm:"size:50;not null"`
	ContactDoctorPhone            string    `json:"contact_doctor_phone" gorm:"size:50;not null"`
	CreatedAt                     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt                     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
