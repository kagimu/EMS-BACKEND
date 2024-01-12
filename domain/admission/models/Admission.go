package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"

	"gorm.io/gorm"
)

type Admission struct {
	gorm.Model
	Number           string               `json:"number" gorm:"size:36;unique;not null"`
	Status           string               `json:"status"`
	AdmissionFeePaid bool                 `json:"admission_fee_paid" gorm:"not null;default:false"`
	StudentNumber    string               `json:"student_number" gorm:"size:36;not null"`
	Student          studentModel.Student `gorm:"foreignKey:StudentNumber;references:APN;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber     string               `json:"campus_number" gorm:"size:36;not null"`
	Campus           campusModels.Campus  `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ClassNumber      string               `json:"class_number" gorm:"size:36;not null"`
	Class            classModels.Class    `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StreamNumber     string               `json:"stream_number" gorm:"size:36;not null"`
	Stream           classModels.Stream   `gorm:"foreignKey:StreamNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber       string               `json:"user_number" gorm:"size:36;not null"`
	User             authModels.User      `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
