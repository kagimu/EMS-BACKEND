package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"time"
)

type MealCourseCount struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Number       string              `json:"number" gorm:"size:36;unique;not null"`
	Count        int                 `json:"count"`
	UserNumber   string              `json:"user_number" gorm:"size:36;not null"`
	User         authModels.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SchoolNumber string              `json:"school_number" gorm:"size:36;not null"`
	School       schoolModels.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
