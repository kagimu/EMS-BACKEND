package requests

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
)

type AssessmentMetricRequest struct {
	MaxScore             float64 `json:"max_score"`
	Description          string  `json:"description"`
	AssessmentAreaNumber string  `json:"assessment_area_number"`
	StudyPeriodNumber    string  `json:"study_period_number"`
	ScheduledOn          string  `json:"scheduled_on"`
	MetricNumber         string  `json:"metric_number"`
	ClassNumber          string  `json:"class_number"`
	CampusNumber         string  `json:"campus_number"`
	RequireComment       string  `json:"require_comment"`
}

type UpdateAssessmentMetricRequest struct {
	AssessmentMetricRequest
	AssessmentMetricNumber string `json:"assessment_metric_number"`
}

func (assessmentMetric *AssessmentMetricRequest) Validate() error {
	if assessmentMetric.AssessmentAreaNumber == "" {
		return errors.New("assessment area number is required")
	}
	if assessmentMetric.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}
	if assessmentMetric.MetricNumber == "" {
		return errors.New("metric number is required")
	}
	if assessmentMetric.CampusNumber == "" {
		return errors.New("Choose campus ")
	}
	if assessmentMetric.ClassNumber == "" {
		return errors.New("Choose a class ")
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

func (assessmentMetric *UpdateAssessmentMetricRequest) Validate() error {

	if assessmentMetric.AssessmentMetricNumber == "" {
		return errors.New("choose the assessment metric being edited")
	}

	err := assessmentMetric.AssessmentMetricRequest.Validate()
	if err != nil {
		return err
	}

	return nil
}
