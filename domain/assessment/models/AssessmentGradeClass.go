package models

import (
	assessmentGradeModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
)

type AssessmentGradeClass struct {
	ID                    uint                                  `gorm:"primarykey"`
	AssessmentGradeNumber string                                `json:"assessment_grade_number" gorm:"size:36;not null"`
	AssessmentGrade       assessmentGradeModels.AssessmentGrade `gorm:"foreignKey:AssessmentGradeNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ClassNumber           string                                `json:"class_number" gorm:"size:36;not null"`
	Class                 classModels.Class                     `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
