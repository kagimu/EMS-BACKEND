package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

type PreAdmissionInterviewTestReport struct {
	gorm.Model
	Number                          string             `json:"number" gorm:"size:36;unique;not null"`
	ApplicationNumber               string             `json:"application_number" gorm:"size:36;not null"`
	Application                     models.Application `gorm:"foreignKey:ApplicationNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AdmissionStatus                 string             `json:"admission_status" gorm:"not null"`
	ClassNumber                     string             `json:"class_number" gorm:"size:36;not null"`
	Class                           classModel.Class   `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AttentionNeeded                 string             `json:"attention_needed" gorm:"type:text;"`
	ReasonForChangingAdmissionClass string             `json:"reason_for_changing_admission_class" gorm:"type:text;"`
	ReasonForNotAdmittingStudent    string             `json:"reason_for_not_admitting_student" gorm:"type:text;"`
	UserNumber                      string             `json:"user_number" gorm:"size:36;not null"`
	User                            authModel.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
