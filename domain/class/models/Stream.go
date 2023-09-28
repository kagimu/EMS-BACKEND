package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Stream struct {
	gorm.Model
	Number      string      `json:"number" gorm:"size:32;not null;unique"`
	Name        string      `json:"name" gorm:"size:100;not null"`
	UserNumber  string      `json:"user_number" gorm:"size:40;not null"`
	User        models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ClassNumber string      `json:"class_number" gorm:"size:40;not null"`
	Class       Class       `gorm:"foreignKey:ClassNumber;references:Number;constraint: onUpdate:CASCADE, onDelete:CASCADE"`
}
