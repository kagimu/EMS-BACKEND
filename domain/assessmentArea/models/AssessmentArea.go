package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type AssessmentArea struct {
	gorm.Model
	Number      string      `json:"number" gorm:"size:36;unique;not null"`
	Name        string      `json:"name" gorm:"not null"`
	ClassNumber string      `json:"class_number" gorm:"size:36;not null"`
	Class       string      `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber  string      `json:"user_number" gorm:"size:36;not null"`
	User        models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
