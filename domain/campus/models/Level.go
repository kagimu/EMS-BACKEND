package models

import (
	levelModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"gorm.io/gorm"
)

type (
	Level struct {
		gorm.Model
		Number              string                   `json:"number" gorm:"size:36;unique;not null"`
		Name                string                   `json:"name" gorm:"-"`
		LearningLevelNumber string                   `json:"learning_level_number" gorm:"size:36;not null"`
		CampusNumber        string                   `json:"campus_number" gorm:"size:36;not null"`
		LearningLevel       levelModel.LearningLevel `gorm:"foreignKey:LearningLevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
	LevelDataInfo struct {
		Name   string `json:"name"`
		Code   int64  `json:"code"`
		Number string `json:"number"`
	}
)

// append name on Level object from LearningLevel via after find hook
func (level *Level) AfterFind(tx *gorm.DB) (err error) {
	var learningLevel levelModel.LearningLevel
	tx.First(&learningLevel, "number = ?", level.LearningLevelNumber)
	level.Name = learningLevel.Name
	return
}
