package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type PreAdmissionInterviewTestReport struct {
	gorm.Model
	Number            string             `json:"number" gorm:"size:36;unique;not null"`
	ApplicationNumber string             `json:"application_number" gorm:"size:36;not null"`
	Application       models.Application `gorm:"foreignKey:ApplicationNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AdmissionStatus   string             `json:"admission_status" gorm:"not null"`
	UserNumber        string             `json:"user_number" gorm:"size:36;not null"`
	User              authModel.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
