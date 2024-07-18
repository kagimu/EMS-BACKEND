package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"time"
)

type UnitOfMeasure struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Number       string              `json:"number" gorm:"primarykey;size:36;not null"`
	Name         string              `json:"name" gorm:"size:50;"`
	Abbreviation string              `json:"abbreviation" gorm:"size:10;"`
	Code         uint64              `json:"code"`
	Status       string              `json:"status"`
	UserNumber   string              `json:"user_number" gorm:"size:36;index;not null"`
	User         authModels.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SchoolNumber string              `json:"school_number" gorm:"size:36;index;not null"`
	School       schoolModels.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
