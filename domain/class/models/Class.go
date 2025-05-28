package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	levelModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"gorm.io/gorm"
)

type (
	Class struct {
		gorm.Model
		ID           uint             `gorm:"primaryKey"`
		Number       string           `json:"number" gorm:"size:36;unique;not null"`
		Name         string           `json:"name" form:"name" binding:"required" gorm:"size:100;not null"`
		Code         int64            `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		Abbreviation string           `json:"abbreviation" form:"abbreviation" binding:"required" gorm:"size:60;not null"`
		UserNumber   string           `json:"user_number" gorm:"size:36;not null"`
		LevelNumber  string           `json:"level_number" gorm:"size:36;not null"`
		Level        levelModel.Level `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User         models.User      `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		EntityNumber string           `json:"entity_number" gorm:"size:36"`
	}

	ClassInfo struct {
		Number   string         `json:"number"`
		Name     string         `json:"name"`
		Subjects []SubjectsInfo `json:"subjects"`
	}

	SubjectsInfo struct {
		Number string `json:"number"`
		Name   string `json:"name"`
	}
)

func (u *Class) TableName() string {
	return "classes"
}
