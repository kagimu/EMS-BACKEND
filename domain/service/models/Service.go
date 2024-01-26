package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Number      string      `json:"number" gorm:"size:36;unique;not null"`
	Description string      `json:"description" gorm:"type:text"`
	Cost        int         `json:"cost"`
	UserNumber  string      `json:"user_number" gorm:"size:36;not null"`
	User        models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

//Wrok on ranking the students by
// Model an invoice with the services paid for. INoice total can be derived from the multiple services
// When the grades are present, someone clicks admit or bulk admit
// Rank based on scores
// Interview number, ask for the number, and admit depedning on the number input.

//ADN gets created if all the services are paid for
