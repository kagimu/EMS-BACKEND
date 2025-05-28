package models

import "gorm.io/gorm"

type EmergencyInformation struct {
	gorm.Model
	StudentNumber                 string `json:"student_number" gorm:"size:36;not null"`
	ContactPerson                 string `json:"contact_person" gorm:"size:50;not null"`
	ContactPersonPhone            string `json:"contact_person_phone" gorm:"size:50;not null"`
	ContactPersonAlternative      string `json:"contact_person_alternative" gorm:"size:50;not null"`
	ContactPersonAlternativePhone string `json:"contact_person_alternative_phone" gorm:"size:50;not null"`
	ContactDoctorName             string `json:"contact_doctor_name" gorm:"size:50;not null"`
	ContactDoctorPhone            string `json:"contact_doctor_phone" gorm:"size:50;not null"`
}
