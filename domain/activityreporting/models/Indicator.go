package models

import (
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"time"
)

type Indicator struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Number       string         `gorm:"unique;type:varchar(255);"`
	Name         string         `gorm:"type:varchar(255);not null"`
	Description  string         `gorm:"type:varchar(1500);not null"`
	Status       string         `gorm:"type:varchar(36);not null"`
	UserNumber   string         `gorm:"type:varchar(36);not null"`
	User         userModel.User `gorm:"foreignkey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	EntityNumber string         `gorm:"type:varchar(36);not null"`
	Entity       models.Entity  `gorm:"foreignkey:EntityNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Code         int            `gorm:"size:10;not null"`
	InputType    string         `gorm:"type:varchar(36);not null"`
}
