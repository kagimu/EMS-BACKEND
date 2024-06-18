package requests

import "errors"

type AggregateScaleRequest struct {
	Code                  int64   `json:"code"`
	MinValue              float64 `json:"min_value"`
	MaxValue              float64 `json:"max_value"`
	CampusNumber          string  `json:"campus_number"`
	GradeNumber           string  `json:"grade_number"`
	Illustration          string  `json:"illustration"`
	GradingStandardNumber string  `json:"grading_standard_number"`
	CampusClassNumber     string  `json:"campus_class_number"`
}

type EditAggregateScaleRequest struct {
	Code                  int64   `json:"code"`
	MinValue              float64 `json:"min_value"`
	MaxValue              float64 `json:"max_value"`
	CampusNumber          string  `json:"campus_number"`
	Illustration          string  `json:"illustration"`
	GradingStandardNumber string  `json:"grading_standard_number"`
	CampusClassNumber     string  `json:"campus_class_number"`
}

func (a *AggregateScaleRequest) Validate() error {

	if a.GradeNumber == "" {
		return errors.New("please select grade  eg D1 , D2 ")
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
	if a.CampusClassNumber == "" {
		return errors.New("at-least one  class number is required")
	}
	if a.CampusNumber == "" {
		return errors.New("campus number is required")
	}
	if a.GradingStandardNumber == "" {
		return errors.New("grading standard number is required")

	}

	return nil
}
