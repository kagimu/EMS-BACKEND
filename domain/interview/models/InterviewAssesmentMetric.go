package models

import (
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"gorm.io/gorm"
)

type InterviewAssessmentMetric struct {
	gorm.Model
	Number                 string                           `json:"number" gorm:"size:36;unique;not null"`
	AssessmentMetricNumber string                           `json:"assessment_metric_number" gorm:"size:36;not null"`
	InterviewNumber        string                           `json:"interview_number" gorm:"size:36;not null"`
	UserNumber             string                           `json:"user_number" gorm:"size:36;not null"`
	Interview              Interview                        `gorm:"foreignKey:InterviewNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentMetric       assessmentModel.AssessmentMetric `gorm:"foreignKey:AssessmentMetricNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
