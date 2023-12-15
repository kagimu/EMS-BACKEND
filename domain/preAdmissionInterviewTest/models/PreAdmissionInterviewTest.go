package models

import "gorm.io/gorm"

type PreAdmissionInterviewTest struct {
	gorm.Model
	Number     string `json:"number" gorm:"size:36;unique;not null"`
	StudentAPN string `json:"student_apn" gorm:"size:36;not null"`
}
