package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	parentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	previousSchoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchool/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type Application struct {
	gorm.Model
	Number                    string                               `json:"number" gorm:"size:36;unique;not null"`
	Status                    string                               `gorm:"type:varchar(255);not null;default:'DRAFT'" json:"status"`
	ApplicationFeePayment     bool                                 `json:"application_fee_payment" gorm:"not null;default:false"`
	ApplicationFeeInvalidDate string                               `json:"application_fee_invalid_date" gorm:"type:varchar(255);"`
	UserNumber                string                               `json:"user_number" gorm:"size:36;index;not null"`
	Student                   studentModel.Student                 `gorm:"foreignKey:ApplicationNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Parent                    parentModel.Parent                   `gorm:"foreignKey:ApplicationNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	PreviousSchools           []previousSchoolModel.PreviousSchool `gorm:"foreignKey:ApplicationNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	User                      models.User                          `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
