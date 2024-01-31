package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Number       string      `json:"number" gorm:"size:36;unique;not null"`
	Description  string      `json:"description" gorm:"type:text"`
	Cost         int         `json:"cost"`
	Name         string      `json:"name"`
	CampusNumber string      `json:"campus_number" gorm:"size:36;not null"`
	UserNumber   string      `json:"user_number" gorm:"size:36;not null"`
	User         models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
