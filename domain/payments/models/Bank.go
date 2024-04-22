package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"gorm.io/gorm"
)

type Bank struct {
	gorm.Model
	Number       string              `json:"number" gorm:"size:36;unique;not null"`
	Name         string              `json:"name" gorm:"size:50;"`
	CampusNumber string              `json:"campus_number" gorm:"size:36;index;not null"`
	Campus       campusModels.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber   string              `json:"user_number" gorm:"size:36;index;not null"`
	User         models.User         `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
