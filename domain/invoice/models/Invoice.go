package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	Number           string         `json:"number" gorm:"size:36;unique;not null"`
	Status           string         `json:"status" gorm:"default:NOT_PAID"`
	StudentNumber    string         `json:"student_number" gorm:"size:36;not null"`
	ServiceNumber    string         `json:"service_number" gorm:"size:36;not null"`
	PaymentReference string         `json:"payment_reference"`
	Student          models.Student `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber       string         `json:"user_number" gorm:"size:36;not null"`
	User             authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
