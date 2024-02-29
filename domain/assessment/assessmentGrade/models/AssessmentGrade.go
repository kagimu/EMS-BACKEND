package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	levelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"

	"gorm.io/gorm"
)

type AssessmentGrade struct {
	gorm.Model
	Number         string            `json:"number" gorm:"size:36;unique;not null"`
	Grade          string            `json:"Grade" gorm:"size:20;"`
	Description    string            `json:"description" gorm:"size:20;"`
	AggregateRange string            `json:"aggregate_range" gorm:"size:40;"`
	Abbreviation   string            `json:"abbreviation" gorm:"size:20;"`
	Condition      string            `json:"condition" gorm:"type:text;"`
	Remarks        string            `json:"remarks" gorm:"type:text;"`
	LevelNumber    string            `json:"level_number" gorm:"size:36;"`
	Level          levelModels.Level `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber     string            `json:"user_number" gorm:"size:36;index;not null"`
	User           authModel.User    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
