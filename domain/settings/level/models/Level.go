package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"time"
)

type (
	Level struct {
		CreatedAt    time.Time
		UpdatedAt    time.Time
		Number       string      `json:"number" gorm:"size:36;unique;not null;primarykey"`
		Name         string      `json:"name" form:"name" binding:"required" gorm:"size:100;not null"`
		Code         int64       `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		UserNumber   string      `json:"user_number" gorm:"size:36;not null"`
		User         models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		EntityNumber string      `json:"entity_number" gorm:"size:36"`
		//Entity       eosModels.Entity `gorm:"foreignKey:EntityNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
	LevelDataInfo struct {
		Name   string `json:"name"`
		Code   int64  `json:"code"`
		Number string `json:"number"`
	}
)

// append name on Level object from LearningLevel via after find hook
/*func (level *Level) AfterFind(tx *gorm.DB) (err error) {
	var learningLevel levelModel.LearningLevel
	tx.First(&learningLevel, "number = ?", level.LearningLevelNumber)
	level.Name = learningLevel.Name
	return
}*/
