package models

import (
	"gorm.io/gorm"
)

type JobInterview struct {
	gorm.Model
	Number           string     `json:"number" gorm:"size:36;unique;index;not null"`
	JobVacancyNumber string     `json:"job_vacancy_number" gorm:"size:36;not null"`
	JobVacancy       JobVacancy `gorm:"foreignKey:JobVacancyNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	InterviewTitle   string     `json:"interview_title" gorm:"size:255;not null"`
	InterviewDate    string     `json:"interview_date" gorm:"size:50"`
	InterviewTime    string     `json:"interview_time" gorm:"size:50"`
	LocationOrLink   string     `json:"location_or_link" gorm:"size:255"`
	Status           string     `json:"status" gorm:"size:50;default:'Scheduled'"`
	UserNumber       string     `json:"user_number" gorm:"size:36;not null"`
}
