package requests

import "errors"

type AggregateScaleRequest struct {
	Code                  int64    `json:"code"`
	MinValue              float64  `json:"min_value"`
	MaxValue              float64  `json:"max_value"`
	Weight                float64  `json:"weight"`
	Description           string   `json:"description"`
	CampusNumber          string   `json:"campus_number"`
	Classes               []string `json:"classes"`
	Name                  string   `json:"name"`
	Illustration          string   `json:"illustration"`
	GradingStandardNumber string   `json:"grading_standard_number"`
	CampusClassNumber     string   `json:"campus_class_number"`
}

func (a *AggregateScaleRequest) Validate() error {

	if a.Name == "" {
		return errors.New("name is required")
	}
	if a.Code == 0 {
		return errors.New("code is required")
	}
	if a.MaxValue < a.MinValue {
		return errors.New("maximum value must be greater than minimum value")
	}
	if a.MaxValue < 0 {
		return errors.New("maximum value must be greater than 0")
	}
	if len(a.Classes) < 1 {
		return errors.New("at-least one  class number is required")
	}
	if a.CampusNumber == "" {
		return errors.New("campus number is required")
	}
	if a.Weight == 0 {
		return errors.New("weight for this scale are required")
	}
	if a.GradingStandardNumber == "" {
		return errors.New("grading standard number is required")

	}

	return nil
}
