package requests

import "errors"

type AssessmentGradeRequest struct {
	Grade             string `json:"grade"`
	Description       string `json:"description"`
	MinimumPoints     int64  `json:"minimum_points"`
	MaximumPoints     int64  `json:"maximum_points"`
	Name              string `json:"name"`
	Abbreviation      string `json:"abbreviation"`
	Remarks           string `json:"remarks"`
	CampusClassNumber string `json:"campus_class_number"`
}

func (r *AssessmentGradeRequest) Validate() error {
	if r.Grade == "" {
		return errors.New("grade is required")
	}
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
	if r.CampusClassNumber == "" {
		return errors.New("campus class number is required")
	}
	return nil
}
