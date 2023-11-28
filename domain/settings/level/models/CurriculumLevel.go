package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	curriculumModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/controllers/models"
	"gorm.io/gorm"
)

type (
	CurriculumLevel struct {
		gorm.Model
		Number           string                     `json:"number" gorm:"size:36;unique;not null"`
		Code             int64                      `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		LevelNumber      string                     `json:"level_number" gorm:"size:36;not null"`
		CurriculumNumber string                     `json:"curriculum_number" gorm:"size:36;not null"`
		UserNumber       string                     `json:"user_number" gorm:"size:36;not null"`
		Level            Level                      `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Curriculum       curriculumModel.Curriculum `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User             models.User                `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

// append name on Level object from LearningLevel via after find hook
/*func (level *Level) AfterFind(tx *gorm.DB) (err error) {
	var learningLevel levelModel.LearningLevel
	tx.First(&learningLevel, "number = ?", level.LearningLevelNumber)
	level.Name = learningLevel.Name
	return
}*/
