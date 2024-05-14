package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"gorm.io/gorm"
)

type AdmissionCondition struct {
	gorm.Model
	AdmissionNumber string          `json:"admission_offer_number" gorm:"size:36"`
	Admission       Admission       `gorm:"foreignKey:AdmissionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ItemNumber      string          `json:"item_number" gorm:"size:36"`
	Item            itemModels.Item `gorm:"foreignKey:ItemNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber      string          `json:"user_number" gorm:"size:36;not null"`
	User            authModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
