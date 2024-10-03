package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type FeedingSession struct {
	gorm.Model
	Number       string             `json:"number" gorm:"size:36;unique;not null"`
	Name         string             `json:"name" gorm:"size:100;not null"`
	Status       string             `json:"status" gorm:"size:50;not null"`
	SchoolNumber string             `json:"school_number" gorm:"size:36;not null"`
	School       schoolModel.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber   string             `json:"user_number" gorm:"size:36;index;not null"`
	User         authModels.User    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
