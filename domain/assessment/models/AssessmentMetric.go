package models

import (
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	campusClassModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

type AssessmentMetric struct {
	gorm.Model
	Number                string                             `json:"number" gorm:"size:36;unique;not null"`
	Weight                string                             `json:"weight"`
	Description           string                             `json:"description" gorm:"type:text;"`
	MaxScore              uint8                              `json:"max_score"`
	ScheduledOn           time.Time                          `json:"scheduled_on" gorm:"size:30"`
	GradingStandardNumber string                             `json:"grading_standard_number"`
	GradingStandard       models.GradingStandard             `gorm:"foreignKey:GradingStandardNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentAreaNumber  string                             `json:"assessment_area_number"`
	AssessmentArea        assessmentAreaModel.AssessmentArea `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;" json:"assessment_area"`
	MetricNumber          string                             `json:"metric_number" gorm:"size:36;"`
	Metric                Metric                             `gorm:"foreignKey:MetricNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;" json:"metric"`
	StudyPeriodNumber     string                             `json:"study_period_number" gorm:"size:36;"`
	StudyPeriod           studyPeriodModel.StudyPeriod       `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusClassNumber     string                             `json:"campus_class_number" gorm:"size:36;"`
	CampusClass           campusClassModel.CampusClass       `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentAreaName    string                             `json:"assessment_area_name" gorm:"-"`
	Name                  string                             `json:"name" gorm:"-"`
}

type AssessmentSubMetric struct {
	gorm.Model
	Number                 string           `json:"number" gorm:"size:36;unique;not null"`
	Name                   string           `json:"name" gorm:"size:255;"`
	AssessmentMetricNumber string           `json:"assessment_metric_number" gorm:"size:36;not null"`
	AssessmentMetric       AssessmentMetric `gorm:"foreignKey:AssessmentMetricNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

// Set value for assessmentAreaname field in afterfind hook
func (assessmentMetric *AssessmentMetric) AfterFind(tx *gorm.DB) (err error) {
	var assessmentArea assessmentAreaModel.AssessmentArea
	if err := tx.First(&assessmentArea, "number = ?", assessmentMetric.AssessmentAreaNumber).Error; err != nil {
		return err
	}

	var metric Metric

	if err := tx.First(&metric, "number = ?", assessmentMetric.MetricNumber).Error; err != nil {
		return err
	}

	assessmentMetric.AssessmentAreaName = assessmentArea.Name
	assessmentMetric.Name = metric.Name + " - " + assessmentArea.Name + " - " + assessmentMetric.ScheduledOn.Format("2006-01-02")
	return nil
}
