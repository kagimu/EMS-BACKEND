package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type AssessmentGrade struct {
	gorm.Model
	Code         int64              `json:"code" gorm:"size:36;not null"`
	Number       string             `json:"number" gorm:"size:36;unique;not null"`
	SchoolNumber string             `json:"school_number" gorm:"size:36;not null"`
	School       schoolModel.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Name         string             `json:"name"`
	Abbreviation string             `json:"abbreviation" gorm:"size:20;"`
	UserNumber   string             `json:"user_number" gorm:"size:36;index;not null"`
	User         authModel.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
