package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

type Billable struct {
	gorm.Model
	Number            string                       `json:"number" gorm:"size:36;unique;not null"`
	UserNumber        string                       `json:"user_number" gorm:"size:36;not null"`
	User              authModel.User               `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriodNumber string                       `json:"study_period_number" gorm:"size:36;not null"`
	StudyPeriod       studyPeriodModel.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ClassNumber       string                       `json:"class_number" gorm:"size:36;not null"`
	Class             classModels.Class            `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	CampusNumber      string                       `json:"campus_number" gorm:"size:36;not null"`
	Campus            campusModel.Campus           `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
