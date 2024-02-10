package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type AdmissionOfferCondition struct {
	gorm.Model
	Number      string          `json:"number" gorm:"size:36;unique;not null"`
	Description string          `json:"description" gorm:"type:text"`
	UserNumber  string          `json:"user_number" gorm:"size:36;not null"`
	User        authModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
