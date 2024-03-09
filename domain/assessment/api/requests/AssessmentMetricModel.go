package requests

import (
	"errors"
	"time"
)

type AssessmentMetricRequest struct {
	Name                  string    `json:"name"`
	Weight                uint8     `json:"weight"`
	MaxScore              uint8     `json:"max_score"`
	ScheduledOn           time.Time `json:"scheduled_on"`
	GradingStandardNumber string    `json:"grading_standard_number"`
	Description           string    `json:"description"`
	AssessmentAreaNumber  string    `json:"assessment_area_number"`
	CurriculumNumber      string    `json:"curriculum_number"`
}

func (a *AssessmentMetricRequest) Validate() error {
	if a.Name == "" {
		return errors.New("name is required")
	}
	if a.Weight == 0 {
		return errors.New("weight is required")
	}
	if a.GradingStandardNumber == "" {
		return errors.New("grading standard number is required")
	}
	if a.ScheduledOn == (time.Time{}) {
		return errors.New("scheduled on date is required")
	}
	if a.AssessmentAreaNumber == "" {
		return errors.New("assessment area number is required")
	}
	if a.CurriculumNumber == "" {
		return errors.New("curriculum number is required")
	}

	return nil
}
