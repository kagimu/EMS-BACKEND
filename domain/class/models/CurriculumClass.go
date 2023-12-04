package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"gorm.io/gorm"
)

type (
	CurriculumClass struct {
		gorm.Model
		Number                string                  `json:"number" gorm:"size:36;unique;not null"`
		ClassNumber           string                  `json:"class_number" gorm:"size:36;not null"`
		CurriculumLevelNumber string                  `json:"curriculum_level_number" gorm:"size:36;not null"`
		Code                  int64                   `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		UserNumber            string                  `json:"user_number" gorm:"size:36;not null"`
		Class                 Class                   `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		CurriculumLevel       models3.CurriculumLevel `gorm:"foreignKey:CurriculumLevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User                  models.User             `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
