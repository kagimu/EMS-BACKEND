package models

import (
	"gorm.io/gorm"
)

type JobApplication struct {
	gorm.Model
	Number            string     `json:"number" gorm:"size:36;unique;index;not null"`
	ApplicationNumber string     `json:"application_number" gorm:"size:100;unique;not null"`
	JobVacancyNumber  string     `json:"job_vacancy_number" gorm:"size:36;not null"`
	JobVacancy        JobVacancy `gorm:"foreignKey:JobVacancyNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	FirstName         string     `json:"first_name" gorm:"size:100;not null"`
	LastName          string     `json:"last_name" gorm:"size:100;not null"`
	Email             string     `json:"email" gorm:"size:150;not null"`
	PhoneNumber       string     `json:"phone_number" gorm:"size:50;not null"`
	Status            string     `json:"status" gorm:"size:50;default:'Pending'"`
	UserNumber        string     `json:"user_number" gorm:"size:36"` // nullable if guest applied, or auth user
}
