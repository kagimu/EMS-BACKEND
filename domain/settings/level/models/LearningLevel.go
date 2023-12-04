package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	curriculumModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	"gorm.io/gorm"
)

type (
	LearningLevel struct {
		gorm.Model
		Number           string                     `json:"number" gorm:"size:36;unique;not null"`
		Name             string                     `json:"name" form:"name" binding:"required" gorm:"size:100;unique;not null"`
		Code             int64                      `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		CurriculumNumber string                     `json:"curriculum_number" gorm:"size:36;not null"`
		UserNumber       string                     `json:"user_number" gorm:"size:36;not null"`
		Curriculum       curriculumModel.Curriculum `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User             models.User                `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
