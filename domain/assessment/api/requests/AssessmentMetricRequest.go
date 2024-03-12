package requests

import (
	"errors"
)

type AssessmentMetricRequest struct {
	Weight                uint8  `json:"weight"`
	MaxScore              uint8  `json:"max_score"`
	GradingStandardNumber string `json:"grading_standard_number"`
	Description           string `json:"description"`
	AssessmentAreaNumber  string `json:"assessment_area_number"`
	StudyPeriodNumber     string `json:"study_period_number"`
}

func (a *AssessmentMetricRequest) Validate() error {
	if a.Weight == 0 {
		return errors.New("weight is required")
	}
	if a.GradingStandardNumber == "" {
		return errors.New("grading standard number is required")
	}
	if a.AssessmentAreaNumber == "" {
		return errors.New("assessment area number is required")
	}
	if a.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}

	return nil
}
