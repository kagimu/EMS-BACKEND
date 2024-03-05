package requests

import "errors"

type AggregateScaleRequest struct {
	MinValue              int    `json:"min_value"`
	MaxValue              int    `json:"max_value"`
	Scale                 string `json:"scale"`
	Description           string `json:"description"`
	StudyPeriodNumber     string `json:"study_period_number"`
	ClassNumber           string `json:"class_number"`
	CampusNumber          string `json:"campus_number"`
	SchoolNumber          string `json:"school_number"`
	Name                  string `json:"name"`
	GradingStandardNumber string `json:"grading_standard_number"`
}

func (a *AggregateScaleRequest) Validate() error {
	if a.Name == "" {
		return errors.New("name is required")
	}
	if a.MaxValue < a.MinValue {
		return errors.New("maximum value must be greater than minimum value")
	}
	if a.MaxValue < 0 {
		return errors.New("maximum value must be greater than 0")
	}
	if a.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}
	if a.ClassNumber == "" {
		return errors.New("class number is required")
	}
	if a.CampusNumber == "" {
		return errors.New("campus number is required")
	}
	if a.Scale == "" {
		return errors.New("scale is required")
	}
	if a.SchoolNumber == "" {
		return errors.New("school number is required")
	}

	if a.GradingStandardNumber == "" {
		return errors.New("grading standard number is required")

	}

	return nil
}
