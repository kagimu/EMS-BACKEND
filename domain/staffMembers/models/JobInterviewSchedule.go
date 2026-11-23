package models

import (
	"gorm.io/gorm"
)

type JobInterviewSchedule struct {
	gorm.Model
	Number               string         `json:"number" gorm:"size:36;unique;index;not null"`
	JobApplicationNumber string         `json:"job_application_number" gorm:"size:36;not null"`
	JobApplication       JobApplication `gorm:"foreignKey:JobApplicationNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	JobInterviewNumber   string         `json:"job_interview_number" gorm:"size:36;not null"`
	JobInterview         JobInterview   `gorm:"foreignKey:JobInterviewNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Remarks              string         `json:"remarks" gorm:"type:text"`
	Comment              string         `json:"comment" gorm:"type:text"`
	Status               string         `json:"status" gorm:"size:50;default:'Scheduled'"`
	UserNumber           string         `json:"user_number" gorm:"size:36;not null"`
}
