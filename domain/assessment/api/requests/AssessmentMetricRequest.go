package requests

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
)

type AssessmentMetricRequest struct {
	MaxScore              float64 `json:"max_score"`
	GradingStandardNumber string  `json:"grading_standard_number"`
	Description           string  `json:"description"`
	AssessmentAreaNumber  string  `json:"assessment_area_number"`
	StudyPeriodNumber     string  `json:"study_period_number"`
	ScheduledOn           string  `json:"scheduled_on"`
	MetricNumber          string  `json:"metric_number"`
	CampusClassNumber     string  `json:"campus_class_number"`
	RequireComment        string  `json:"require_comment"`
}

func (assessmentMetric *AssessmentMetricRequest) Validate() error {
	if assessmentMetric.GradingStandardNumber == "" {
		return errors.New("grading standard number is required")
	}
	if assessmentMetric.AssessmentAreaNumber == "" {
		return errors.New("assessment area number is required")
	}
	if assessmentMetric.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}
	if assessmentMetric.MetricNumber == "" {
		return errors.New("metric number is required")
	}
	if assessmentMetric.CampusClassNumber == "" {
		return errors.New("campus class number is required")
	}
	if assessmentMetric.ScheduledOn == "" {
		return errors.New("scheduled on is required")
	}
	if assessmentMetric.RequireComment == "" {
		return errors.New("please specify if the comment  is required or not")
	}
	if assessmentMetric.RequireComment != constants.StatusYes && assessmentMetric.RequireComment != constants.StatusNo {
		return errors.New("require comment should be YES or NO")
	}
	if assessmentMetric.MaxScore <= 0 {
		return errors.New("max score can not be zero or negative")
	}

	return nil
}
