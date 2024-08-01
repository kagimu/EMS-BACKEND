package models

import (
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
)

type InterviewAssessmentMetric struct {
	Number                 string                           `json:"number" gorm:"size:36;unique;not null;primaryKey"`
	AssessmentMetricNumber string                           `json:"assessment_metric_number" gorm:"size:36;not null"`
	InterviewNumber        string                           `json:"interview_number" gorm:"size:36;not null"`
	Interview              Interview                        `gorm:"foreignKey:InterviewNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	AssessmentMetric       assessmentModel.AssessmentMetric `gorm:"foreignKey:AssessmentMetricNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
