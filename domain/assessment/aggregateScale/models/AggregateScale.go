package models

import (
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
)

type AggregateScale struct {
	Code                  int64                   `json:"code" gorm:"size:36;not null"`
	Number                string                  `json:"number" gorm:"size:36;unique;not null"`
	MinValue              float64                 `json:"min_value" gorm:"not null"`
	MaxValue              float64                 `json:"max_value" gorm:"not null"`
	Remarks               string                  `json:"remarks" gorm:"size:1000;not null"`
	ScaleNumber           string                  `json:"scale_number" gorm:"size:36;not null;primarykey"`
	Scale                 models2.Scale           `gorm:"foreignKey:ScaleNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Weight                float64                 `json:"weight"`
	Illustration          string                  `json:"illustration" form:"illustration" gorm:"type:text;"`
	AssessmentGradeNumber string                  `json:"assessment_grade_number" gorm:"size:36;not null;primarykey"`
	AssessmentGrade       models2.AssessmentGrade `gorm:"foreignKey:AssessmentGradeNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
