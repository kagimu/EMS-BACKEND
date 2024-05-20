package requests

import "errors"

type AssessmentGradeRequest struct {
	Code              int64   `json:"code"`
	Grade             string  `json:"grade"`
	Description       string  `json:"description"`
	MinimumPoints     float64 `json:"minimum_points"`
	MaximumPoints     float64 `json:"maximum_points"`
	Name              string  `json:"name"`
	Abbreviation      string  `json:"abbreviation"`
	Remarks           string  `json:"remarks"`
	CampusClassNumber string  `json:"campus_class_number"`
	CampusNumber      string  `json:"campus_number"`
}

type EditAssessmentGradeRequest struct {
	Code              int64   `json:"code"`
	Grade             string  `json:"grade"`
	Description       string  `json:"description"`
	MinimumPoints     float64 `json:"minimum_points"`
	MaximumPoints     float64 `json:"maximum_points"`
	Name              string  `json:"name"`
	Abbreviation      string  `json:"abbreviation"`
	Remarks           string  `json:"remarks"`
	CampusClassNumber string  `json:"campus_class_number"`
	CampusNumber      string  `json:"campus_number"`
}

func (r *AssessmentGradeRequest) Validate() error {
	if r.Abbreviation == "" {
		return errors.New("abbreviation is required")
	}
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Description == "" {
		return errors.New("description is required")
	}
	if r.MinimumPoints == 0 {
		return errors.New("minimum points is required")
	}
	if r.MaximumPoints == 0 {
		return errors.New("maximum points is required")
	}

	if r.Remarks == "" {
		return errors.New("remarks is required")
	}
	if r.CampusNumber == "" {
		return errors.New("campus number is required")
	}
	if r.Code == 0 {
		return errors.New("code is required")
	}
	if r.CampusClassNumber == "" {
		return errors.New("at-least one class is required")
	}
	return nil
}
