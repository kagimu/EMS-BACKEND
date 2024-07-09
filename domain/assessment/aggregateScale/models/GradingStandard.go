package models

import (
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)
import userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"

type GradingStandard struct {
	gorm.Model
	Code         int64               `json:"code" gorm:"size:36;not null"`
	Number       string              `json:"number" gorm:"size:36;unique;not null"`
	Name         string              `json:"name" gorm:"size:255;not null"`
	Description  string              `json:"description" gorm:"type:text;"`
	Type         string              `json:"type" gorm:"size:255;not null"`
	SchoolNumber string              `json:"school_number" gorm:"size:36;not null"`
	School       schoolModels.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber   string              `json:"user_number" gorm:"size:36;not null"`
	User         userModel.User      `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
