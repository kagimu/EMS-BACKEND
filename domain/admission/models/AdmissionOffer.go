package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type AdmissionOffer struct {
	gorm.Model
	Number       string `gorm:"size:36;not null"`
	ClassNumber  string
	SchoolNumber string
	CampusNumber string
	StudentAPN   string
	UserNumber   string         `json:"user_number" gorm:"size:36;not null"`
	User         authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
