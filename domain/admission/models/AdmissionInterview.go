package models

import (
	interviewModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	"gorm.io/gorm"
)

type AdmissionInterview struct {
	gorm.Model
	Number          string                   `json:"number" gorm:"size:36;unique;not null;primaryKey"`
	InterviewNumber string                   `json:"interview_number" gorm:"size:36;not null"`
	Interview       interviewModel.Interview `gorm:"foreignKey:InterviewNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AdmissionNumber string                   `json:"admission_number" gorm:"size:36;not null;unique"`
	Admission       Admission                `gorm:"foreignKey:AdmissionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
