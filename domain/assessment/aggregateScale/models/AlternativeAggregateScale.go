package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

type AlternativeAggregateScale struct {
	gorm.Model
	Number            string                       `json:"number" gorm:"size:36;unique;not null"`
	Name              string                       `json:"name" gorm:"type:text"`
	Factors           string                       `json:"factors" gorm:"type:text"`
	Description       string                       `json:"description" gorm:"type:text"`
	MinValue          int                          `json:"min_value" gorm:"not null"`
	MaxValue          int                          `json:"max_value" gorm:"not null"`
	LevelNumber       string                       `json:"level_number" gorm:"size:36;not null"`
	Level             models.Level                 `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriodNumber string                       `json:"study_period_number" gorm:"size:36;not null"`
	StudyPeriod       studyPeriodModel.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ClassNumber       string                       `json:"class_number" gorm:"size:36;not null"`
	Class             classModel.Class             `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber      string                       `json:"campus_number" gorm:"size:36;not null"`
	Campus            campusModel.Campus           `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber        string                       `json:"user_number" gorm:"size:36;not null"`
	User              authModel.User               `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
