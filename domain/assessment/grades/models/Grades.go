package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type Grade struct {
	gorm.Model
	Number       string             `json:"number" gorm:"size:36;not null;unique"`
	Name         string             `json:"name" gorm:"size:36;not null"`
	ShortName    string             `json:"short_name" gorm:"size:36;not null"`
	SchoolNumber string             `json:"school_number" gorm:"size:36;not null"`
	Code         int64              `json:"code" gorm:"size:36;not null"`
	UserNumber   string             `json:"user_number" gorm:"size:36;not null"`
	Status       string             `json:"status" gorm:"size:36;not null;default:ACTIVE"`
	GradeNumber  string             `json:"grade_number" gorm:"size:36;default:4524f5f5-0c65-45cb-8a7f-2238bd3685cc"`
	User         authModel.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	School       schoolModel.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type GradeGroup struct {
	gorm.Model
	Number       string             `json:"number" gorm:"size:36;not null;unique"`
	Name         string             `json:"name" gorm:"size:36;not null"`
	Code         int64              `json:"code" gorm:"size:36;not null"`
	UserNumber   string             `json:"user_number" gorm:"size:36;not null"`
	User         authModel.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SchoolNumber string             `json:"school_number" gorm:"size:36;not null"`
	School       schoolModel.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
