package models

import (
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"time"
)

type ScriptMarkingGuide struct {
	ID                     uint `gorm:"primaryKey"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	ScriptName             string           `gorm:"column:script_name;size:64;unique"`
	MarkingGuideName       string           `gorm:"column:marking_guide_name;size:64;unique"`
	UserNumber             string           `gorm:"size:36;not null;column:user_number"`
	User                   userModel.User   `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentMetricNumber string           `gorm:"size:36;not null;column:assessment_metric_number"`
	AssessmentMetric       AssessmentMetric `gorm:"foreignKey:AssessmentMetricNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
