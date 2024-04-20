package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	curriculumLevelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"gorm.io/gorm"
)

type ScopeUserLevel struct {
	gorm.Model
	Number        string                      `json:"number" gorm:"size:36;unique;not null"`
	UserNumber    string                      `json:"user_number" gorm:"size:36;not null"`
	User          authModel.User              `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	LevelNumber   string                      `json:"level_number" gorm:"size:36;not null"`
	Level         curriculumLevelModels.Level `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AddedByNumber string                      `json:"added_by_number" gorm:"size:36;not null"`
	AddedBy       authModel.User              `gorm:"foreignKey:AddedByNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
