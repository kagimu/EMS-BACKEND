package models

import (
	"gorm.io/gorm"
)

type (
	LessonPlanStep struct {
		gorm.Model
		Number                  string             `json:"number" gorm:"size:36;unique;not null"`
		Code                    int64              `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		Step                    string             `json:"step" form:"name" binding:"required" gorm:"type:longtext;not null"`
		LessonPlanConceptNumber string             `json:"lesson_plan_concept_number" gorm:"size:36;not null"`
		LessonPlanConcepts      LessonPlanConcepts `gorm:"foreignKey:LessonPlanConceptNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
