package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type SalaryScale struct {
	gorm.Model
	Number        string      `json:"number" gorm:"size:36;unique;not null"`
	MinimumAmount uint64      `json:"minimum"`
	MaximumAmount uint64      `json:"maximum"`
	Name          string      `json:"name" gorm:"unique"`
	Description   string      `json:"description" gorm:"type:text"`
	UserNumber    string      `json:"user_number" gorm:"size:36;index;not null"`
	User          models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
