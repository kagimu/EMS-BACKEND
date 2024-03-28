package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusClass "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"

	"gorm.io/gorm"
)

type AssessmentGrade struct {
	gorm.Model
	Number            string                  `json:"number" gorm:"size:36;unique;not null"`
	Grade             string                  `json:"Grade" gorm:"size:20;"`
	Name              string                  `json:"name"`
	Description       string                  `json:"description" gorm:"type:text"`
	MinimumPoints     int64                   `json:"minimum_points"`
	MaximumPoints     int64                   `json:"maximum_points"`
	Abbreviation      string                  `json:"abbreviation" gorm:"size:20;"`
	Remarks           string                  `json:"remarks" gorm:"type:text;"`
	CampusClassNumber string                  `json:"campus_class_number" gorm:"size:36;not null"`
	CampusClass       campusClass.CampusClass `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber        string                  `json:"user_number" gorm:"size:36;index;not null"`
	User              authModel.User          `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
