package models

import (
	"time"

	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	staffModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"gorm.io/gorm"
)

type Feeding struct {
	gorm.Model
	Number               string            `json:"number" gorm:"size:36;unique;not null"`
	FeedingSessionNumber string            `json:"feeding_session_number" gorm:"size:36;"`
	Date                 time.Time         `json:"date" gorm:"not null"`
	StockingNumber       string            `json:"stocking_number" gorm:"size:36;"`
	Stocking             Stocking          `gorm:"foreignKey:StockingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StaffNumber          string            `json:"staff_number" gorm:"size:36;"`
	Staff                staffModels.Staff `gorm:"foreignKey:StaffNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	FeedingSession       FeedingSession    `gorm:"foreignKey:FeedingSessionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber           string            `json:"user_number" gorm:"size:36;index;not null"`
	User                 authModels.User   `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
