package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type PreAdmissionInterviewTest struct {
	gorm.Model
	Number          string         `json:"number" gorm:"size:36;unique;not null"`
	StudentNumber   string         `json:"student_number" gorm:"size:36;not null"`
	Student         models.Student `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AdmissionStatus string         `json:"admission_status"`
}

type Subject struct {
	gorm.Model
	Name                       string `json:"name" gorm:"size:36;not null"`
	PreAdmissionInterviewTests []PreAdmissionInterviewTest
}

func (p *PreAdmissionInterviewTest) TableName() string {
	return "pre_interview_test"
}
