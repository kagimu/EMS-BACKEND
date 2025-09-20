package models

import (
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/models"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
)

type ScriptMarkingGuide struct {
	ID                     uint `gorm:"primaryKey"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	ScriptNumber           string            `gorm:"column:script_number;size:36;unique"`
	Script                 models.Attachment `gorm:"foreignKey:ScriptNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	MarkingGuideNumber     string            `gorm:"column:marking_guide_number;size:36;unique"`
	MarkingGuide           models.Attachment `gorm:"foreignKey:MarkingGuideNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber             string            `gorm:"size:36;not null;column:user_number"`
	User                   userModel.User    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentMetricNumber string            `gorm:"size:36;not null;column:assessment_metric_number"`
	AssessmentMetric       AssessmentMetric  `gorm:"foreignKey:AssessmentMetricNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
