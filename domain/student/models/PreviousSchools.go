package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

type PreviousSchools struct {
	gorm.Model
	Number        string              `json:"number" gorm:"size:36;unique;not null"`
	Name          string              `json:"name" gorm:"size:50;not null"`
	StudentNumber string              `json:"student_number" gorm:"size:36;unique;not null"`
	Classes       []classModels.Class `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber    string              `json:"user_number" gorm:"size:36;index;not null"`
	User          models.User         `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
