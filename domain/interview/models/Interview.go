package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Interview struct {
	gorm.Model
	Number        string         `json:"number" gorm:"size:36;unique;not null"`
	Message       string         `json:"message"`
	Venue         string         `json:"venue" gorm:"size:100;not null"`
	Name          string         `json:"name" gorm:"size:100;"`
	InterviewDate string         `json:"interview_date" gorm:"size:100;not null"`
	UserNumber    string         `json:"user_number" gorm:"size:36;not null"`
	User          authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
