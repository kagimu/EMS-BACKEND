package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type UnitOfMeasure struct {
	gorm.Model
	Number       string          `json:"number" gorm:"size:36;unique;not null"`
	Name         string          `json:"name" gorm:"size:50;"`
	Abbreviation string          `json:"abbreviation" gorm:"size:10;"`
	Code         uint64          `json:"code"`
	UserNumber   string          `json:"user_number" gorm:"size:36;index;not null"`
	User         authModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
