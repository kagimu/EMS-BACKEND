package models

import (
	"gorm.io/gorm"
)

type (
	LessonPlanSteps struct {
		gorm.Model
		Number     string     `json:"number" gorm:"size:36;unique;not null"`
		Code       int64      `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		Step       string     `json:"step" form:"name" binding:"required" gorm:"type:longtext;not null"`
		LessonPlan LessonPlan `gorm:"foreignKey:SchemeSubTopicNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
