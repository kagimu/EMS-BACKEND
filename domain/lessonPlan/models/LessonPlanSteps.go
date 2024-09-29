package models

import (
	"time"
)

type (
	LessonPlanStep struct {
		ID                      uint `gorm:"primaryKey"`
		CreatedAt               time.Time
		UpdatedAt               time.Time
		Number                  string            `json:"number" gorm:"size:36;unique;not null"`
		Code                    int64             `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		Step                    string            `json:"step" form:"name" binding:"required" gorm:"type:longtext;not null"`
		LessonPlanConceptNumber string            `json:"lesson_plan_concept_number" gorm:"size:36;not null"`
		LessonPlanConcept       LessonPlanConcept `gorm:"foreignKey:LessonPlanConceptNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	}
)
