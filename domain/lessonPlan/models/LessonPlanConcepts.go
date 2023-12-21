package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	conceptModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	"gorm.io/gorm"
)

type (
	LessonPlanConcepts struct {
		gorm.Model
		Number           string               `json:"number" gorm:"size:36;unique;not null"`
		Code             int64                `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		ConceptNumber    string               `json:"concept_number" form:"name" binding:"required" gorm:"size:36;not null"`
		LessonPlanNumber string               `json:"lesson_plan_number" form:"name" binding:"required" gorm:"size:36;not null"`
		UserNumber       string               `json:"user_number" gorm:"size:36;not null"`
		LessonPlan       LessonPlan           `gorm:"foreignKey:LessonPlanNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Concept          conceptModel.Concept `gorm:"foreignKey:ConceptNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User             models.User          `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
