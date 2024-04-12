package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type FeedingSessionPerDay struct {
	gorm.Model
	Number               string          `json:"number" gorm:"size:36;unique;not null"`
	FeedingSessionNumber string          `json:"feeding_session_number" gorm:"size:36;"`
	Day                  string          `json:"day" gorm:"not null"`
	FeedingSession       FeedingSession  `gorm:"foreignKey:FeedingSessionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber           string          `json:"user_number" gorm:"size:36;index;not null"`
	User                 authModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
