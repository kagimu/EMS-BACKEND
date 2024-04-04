package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type Assessment struct {
	gorm.Model
	Number                 string               `json:"number" gorm:"size:36;unique;not null"`
	AssessmentMetricNumber string               `json:"assessment_metric_number" gorm:"size:36;"`
	AssessmentMetric       AssessmentMetric     `gorm:"foreignKey:AssessmentMetricNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudentNumber          string               `json:"student_number" gorm:"size:36;index;not null"`
	Student                studentModel.Student `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;" json:"student"`
	Comment                string               `json:"comment" gorm:"size:250"`
	ActualScore            float64              `json:"actual_score"`
	UserNumber             string               `json:"user_number" gorm:"size:36;index;not null"`
	User                   authModel.User       `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
