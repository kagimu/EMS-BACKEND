package models

import (
	"time"

	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
)

type ModeOfStudy struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Number       string         `json:"number" gorm:"size:36;unique;not null"`
	UserNumber   string         `json:"user_number" gorm:"size:36;not null"`
	User         authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Status       string         `json:"status" gorm:"size:64;not null;default:ACTIVE"`
	Name         string         `json:"name" gorm:"size:64;not null"`
	SchoolNumber string         `json:"school_number" gorm:"size:36;not null"`
	School       School         `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Code         int64          `json:"code"`
}

func (modeOfStudy *ModeOfStudy) TableName() string {
	return "mode_of_studies"
}
