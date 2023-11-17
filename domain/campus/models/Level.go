package models

import (
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"gorm.io/gorm"
)

type (
	Level struct {
		gorm.Model
		Number              string                `json:"number" gorm:"size:36;unique;not null"`
		LearningLevelNumber string                `json:"learning_level_number" gorm:"size:36;not null"`
		CampusNumber        string                `json:"campus_number" gorm:"size:36;not null"`
		LearningLevel       models2.LearningLevel `gorm:"foreignKey:LearningLevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
	LevelDataInfo struct {
		Name   string `json:"name"`
		Number string `json:"number"`
	}
)
