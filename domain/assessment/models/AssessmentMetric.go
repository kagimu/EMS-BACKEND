package models

import (
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	"gorm.io/gorm"
)

type AssessmentMetric struct {
	gorm.Model
	Number               string                             `json:"number" gorm:"size:36;unique;not null"`
	Name                 string                             `json:"name" gorm:"size:255;not null"`
	Weight               uint8                              `json:"weight" gorm:"default:1"`
	Description          string                             `json:"description" gorm:"type:text;"`
	MaxScore             uint8                              `json:"max_score"`
	ActualScore          uint8                              `json:"actual_score"`
	Routine              string                             `json:"routine" gorm:"size:30"`
	Schedule             string                             `json:"schedule" gorm:"size:30"`
	AssessmentAreaNumber string                             `json:"assessment_area_number"`
	AssessmentArea       assessmentAreaModel.AssessmentArea `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type AssessmentSubMetric struct {
	gorm.Model
	Number                 string           `json:"number" gorm:"size:36;unique;not null"`
	Name                   string           `json:"name" gorm:"size:255;"`
	AssessmentMetricNumber string           `json:"assessment_metric_number" gorm:"size:36;not null"`
	AssessmentMetric       AssessmentMetric `gorm:"foreignKey:AssessmentMetricNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
