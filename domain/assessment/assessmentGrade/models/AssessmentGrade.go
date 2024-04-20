package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusClass "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"

	"gorm.io/gorm"
)

type AssessmentGrade struct {
	gorm.Model
	Code              int64                   `json:"code" gorm:"size:36;not null"`
	Number            string                  `json:"number" gorm:"size:36;unique;not null"`
	Name              string                  `json:"name"`
	Description       string                  `json:"description" gorm:"type:text"`
	MinimumPoints     float64                 `json:"minimum_points"`
	MaximumPoints     float64                 `json:"maximum_points"`
	Abbreviation      string                  `json:"abbreviation" gorm:"size:20;"`
	Illustration      string                  `json:"illustration" form:"illustration" gorm:"type:text"`
	Remarks           string                  `json:"remarks" gorm:"type:text;"`
	CampusNumber      string                  `json:"campus_number" gorm:"size:36;not null"`
	Campus            campusClass.CampusClass `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusClassNumber string                  `json:"campus_class_number" gorm:"size:36;not null"`
	CampusClass       campusClass.CampusClass `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber        string                  `json:"user_number" gorm:"size:36;index;not null"`
	User              authModel.User          `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
