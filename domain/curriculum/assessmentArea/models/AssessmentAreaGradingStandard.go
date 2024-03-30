package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	"gorm.io/gorm"
)

type AssessmentAreaGradingStandard struct {
	gorm.Model
	Number                string                 `json:"number" gorm:"size:36;unique;not null"`
	AssessmentAreaNumber  string                 `json:"assessment_area_number" gorm:"size:36; not null"`
	GradingStandardNumber string                 `json:"grading_standard_number" gorm:"size:36;not null"`
	GradingStandard       models.GradingStandard `gorm:"foreignKey:GradingStandardNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentArea        AssessmentArea         `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
