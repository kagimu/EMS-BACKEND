package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"time"
)

type LandingPage struct {
	ID          uint   `gorm:"primarykey"`
	Number      string `json:"number" gorm:"size:36;unique;not null"`
	CreatedAt   time.Time
	Title       string      `json:"title" gorm:"size:120;index;not null"`
	Description string      `json:"description"`
	Image       string      `json:"image" gorm:"size:64;index;not null"`
	UserNumber  string      `json:"user_number" gorm:"size:36;index;not null"`
	User        models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
