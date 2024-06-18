package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type FormatsOfMeasure struct {
	gorm.Model
	Number       string             `json:"number" gorm:"size:36;not null;unique"`
	Name         string             `json:"name" gorm:"size:36;not null; primaryKey"`
	Weight       float64            `json:"weight"`
	ShortName    string             `json:"short_name" gorm:"primaryKey"`
	SchoolNumber string             `json:"school_number" gorm:"size:36;not null"`
	Description  string             `json:"description" gorm:"size:255"`
	Code         int64              `json:"code" gorm:"size:36;not null"`
	UserNumber   string             `json:"user_number" gorm:"size:36;not null"`
	User         authModel.User     `gorm:"foreignKey:UserNumber;references:Number"`
	School       schoolModel.School `gorm:"foreignKey:SchoolNumber;references:Number"`
}
