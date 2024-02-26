package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

type AssessmentMetric struct {
	gorm.Model
	Number                    string              `json:"number" gorm:"size:36;unique;not null"`
	Name                      string              `json:"name" gorm:"size:255;not null"`
	Weight                    uint8               `json:"weight" gorm:"default:1"`
	PercentageScore           uint8               `json:"percentage_score" gorm:"max:100"`
	DescriptiveScore          string              `json:"descriptive_score" gorm:"type:text"`
	ExpectedScore             uint8               `json:"expected_score"`
	ActualScore               uint8               `json:"actual_score"`
	ClassNumber               string              `json:"class_number" gorm:"size:36;"`
	Schedule                  string              `json:"schedule" gorm:"size:30"`
	Class                     models.Class        `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentSubMetricNumber string              `json:"assessment_metric_number" gorm:"size:36;not null"`
	AssessmentSubMetric       AssessmentSubMetric `gorm:"foreignKey:AssessmentSubMetricNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type AssessmentSubMetric struct {
	gorm.Model
	Number string `json:"number" gorm:"size:36;unique;not null"`
	Name   string `json:"name" gorm:"size:255;"`
}
