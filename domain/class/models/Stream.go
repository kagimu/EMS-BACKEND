package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Stream struct {
	gorm.Model
	Number            string      `json:"number" gorm:"size:36;not null;unique"`
	Name              string      `json:"name" gorm:"size:200;not null"`
	Status            string      `json:"status" gorm:"size:20;default:'ACTIVE'"`
	UserNumber        string      `json:"user_number" gorm:"size:36;not null"`
	User              models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusClassNumber string      `json:"campus_class_number" gorm:"size:36;not null"`
	CampusClass       CampusClass `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: onUpdate:CASCADE, onDelete:CASCADE"`
}
