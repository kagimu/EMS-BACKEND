package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type Department struct {
	gorm.Model
	Number       string              `json:"number" gorm:"size:36;unique;index;not null"`
	SchoolNumber string              `json:"school_number" gorm:"size:36;not null"`
	School       schoolModels.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Name         string              `json:"name" form:"name" binding:"required" gorm:"size:100;not null"`
	Code         int64               `json:"code" form:"code" binding:"required" gorm:"not null"`
	UserNumber   string              `json:"user_number" gorm:"size:36;not null"`
	User         authModels.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
