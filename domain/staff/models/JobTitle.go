package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type JobTitle struct {
	gorm.Model
	Number      string      `json:"number" gorm:"size:36;unique;not null"`
	Name        string      `json:"name" gorm:"unique"`
	Description string      `json:"description"`
	UserNumber  string      `json:"user_number" gorm:"size:36;index;not null"`
	User        models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
