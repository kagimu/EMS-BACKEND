package models

import (
	userModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"gorm.io/gorm"
)

type Grade struct {
	gorm.Model
	Number      string          `json:"number" gorm:"size:36;unique;not null"`
	LevelNumber string          `json:"level_number" gorm:"size:36;unique;not null"`
	Level       models.Level    `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber  string          `json:"user_number" gorm:"size:36;index;not null"`
	User        userModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
