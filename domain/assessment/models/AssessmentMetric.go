package models

import (
	"time"

	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"

	campusClassModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

type AssessmentMetric struct {
	gorm.Model
	Number               string                             `json:"number" gorm:"size:36;unique;not null"`
	Description          string                             `json:"description" gorm:"type:text;"`
	MaxScore             float64                            `json:"max_score"`
	ScheduledOn          time.Time                          `json:"scheduled_on" gorm:"size:30"`
	AssessmentAreaNumber string                             `json:"assessment_area_number"`
	AssessmentArea       assessmentAreaModel.AssessmentArea `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;" json:"assessment_area"`
	MetricNumber         string                             `json:"metric_number" gorm:"size:36;"`
	Metric               Metric                             `gorm:"foreignKey:MetricNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;" json:"metric"`
	StudyPeriodNumber    string                             `json:"study_period_number" gorm:"size:36;"`
	StudyPeriod          studyPeriodModel.StudyPeriod       `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusClassNumber    string                             `json:"campus_class_number" gorm:"size:36;"`
	CampusClass          campusClassModel.CampusClass       `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	RequireComment       string                             `json:"require_comment" gorm:"size:5;not null;default:'YES'"`
	UserNumber           string                             `json:"user_number" gorm:"size:36;not null"`
	User                 userModel.User                     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;" json:"user"`
	MarkingGuide         string                             `json:"marking_guide" gorm:"size:36;"`
	TestScript           string                             `json:"test_script" gorm:"size:36;"`
}
