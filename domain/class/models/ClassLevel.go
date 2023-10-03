package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"gorm.io/gorm"
)

type (
	ClassLevel struct {
		gorm.Model
		Number       string         `json:"number" gorm:"size:36;unique;not null"`
		ClassNumber  string         `json:"class_number" gorm:"size:36;not null"`
		CampusNumber string         `json:"campus_number" gorm:"size:36;not null"`
		LevelNumber  string         `json:"level_number" gorm:"size:36;not null"`
		UserNumber   string         `json:"user_number" gorm:"size:36;not null"`
		Class        Class          `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Campus       models2.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Level        models2.Level  `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User         models.User    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *ClassLevel) TableName() string {
	return "class_levels"
}
