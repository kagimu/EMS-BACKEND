package requests

import "errors"

type AlternativeAggregateScaleRequest struct {
	Name              string `json:"name"`
	Factors           string `json:"factors"`
	Description       string `json:"description"`
	MinValue          int    `json:"min_value"`
	MaxValue          int    `json:"max_value"`
	Scale             string `json:"scale"`
	LevelNumber       string `json:"level_number"`
	StudyPeriodNumber string `json:"study_period_number"`
	ClassNumber       string `json:"class_number"`
	CampusNumber      string `json:"campus_number"`
}

func (a *AlternativeAggregateScaleRequest) Validate() error {
	if a.Description == "" {
		return errors.New("description is required")
	}
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

	return nil
}
