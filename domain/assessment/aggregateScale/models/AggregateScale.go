package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type AggregateScale struct {
	gorm.Model
	Number                string             `json:"number" gorm:"size:36;unique;not null"`
	MinValue              int                `json:"min_value" gorm:"not null"`
	MaxValue              int                `json:"max_value" gorm:"not null"`
	Name                  string             `json:"name" gorm:"type:text"`
	Points                int64              `json:"points" gorm:"not null"`
	SchoolNumber          string             `json:"school_number" gorm:"size:36;not null"`
	Description           string             `json:"description" gorm:"type:text"`
	School                schoolModel.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ClassNumber           string             `json:"class_number" gorm:"size:36;not null"`
	Class                 classModel.Class   `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber          string             `json:"campus_number" gorm:"size:36;not null"`
	Campus                campusModel.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber            string             `json:"user_number" gorm:"size:36;not null"`
	User                  authModel.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	GradingStandardNumber string             `json:"grading_standard_number" gorm:"size:36;not null"`
	GradingStandard       GradingStandard    `gorm:"foreignKey:GradingStandardNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
