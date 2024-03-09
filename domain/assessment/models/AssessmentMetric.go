package models

import (
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

type AssessmentMetric struct {
	gorm.Model
	Number                string                             `json:"number" gorm:"size:36;unique;not null"`
	Name                  string                             `json:"name" gorm:"size:255;not null"`
	Weight                uint8                              `json:"weight" gorm:"default:1"`
	Description           string                             `json:"description" gorm:"type:text;"`
	MaxScore              uint8                              `json:"max_score"`
	ScheduledOn           time.Time                          `json:"scheduled_on" gorm:"size:30"`
	GradingStandardNumber string                             `json:"grading_standard_number"`
	GradingStandard       models.GradingStandard             `gorm:"foreignKey:GradingStandardNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentAreaNumber  string                             `json:"assessment_area_number"`
	AssessmentArea        assessmentAreaModel.AssessmentArea `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriodNumber     string                             `json:"study_period_number" gorm:"size:36;"`
	StudyPeriod           studyPeriodModel.StudyPeriod       `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type AssessmentSubMetric struct {
	gorm.Model
	Number                 string           `json:"number" gorm:"size:36;unique;not null"`
	Name                   string           `json:"name" gorm:"size:255;"`
	AssessmentMetricNumber string           `json:"assessment_metric_number" gorm:"size:36;not null"`
	AssessmentMetric       AssessmentMetric `gorm:"foreignKey:AssessmentMetricNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
