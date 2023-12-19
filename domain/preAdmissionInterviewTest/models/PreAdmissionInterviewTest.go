package models

import "gorm.io/gorm"

type PreAdmissionInterviewTest struct {
	gorm.Model
	Number     string `json:"number" gorm:"size:36;unique;not null"`
	StudentAPN string `json:"student_apn" gorm:"size:36;not null"`
	Scores     string
}

type Subject struct {
	gorm.Model
	Name                       string `json:"name" gorm:"size:36;not null"`
	PreAdmissionInterviewTests []PreAdmissionInterviewTest
}

func (p *PreAdmissionInterviewTest) TableName() string {
	return "pre_interview_test"
}
