package requests

import "errors"

type AssessmentGradeRequest struct {
	Classes              []string                     `json:"classes"`
	AssessmentAreaCount  int                          `json:"assessment_area_count"`
	AssessmentGradeItems []AssessmentGradeItemRequest `json:"assessment_grade_items"`
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

type AssessmentGradeItemRequest struct {
	Code          int64   `json:"code"`
	Description   string  `json:"description"`
	MinimumPoints float64 `json:"minimum_points"`
	MaximumPoints float64 `json:"maximum_points"`
	Name          string  `json:"name"`
	Abbreviation  string  `json:"abbreviation"`
	Remarks       string  `json:"remarks"`
}

func (assessmentRequest *AssessmentGradeRequest) Validate() error {

	if len(assessmentRequest.Classes) <= 0 {
		return errors.New("choose classes")
	}

	if assessmentRequest.AssessmentAreaCount <= 0 {
		return errors.New("set the number of assessment areas to consider I.E best 8")
	}

	for _, r := range assessmentRequest.AssessmentGradeItems {

		if r.Abbreviation == "" {
			return errors.New("abbreviation is required")
		}

		if r.Name == "" {
			return errors.New("name is required")
		}

		if r.Description == "" {
			return errors.New("description is required")
		}

		if r.MinimumPoints <= 0 {
			return errors.New("minimum points is required")
		}

		if r.MaximumPoints <= 0 {
			return errors.New("maximum points is required")
		}

		if r.Remarks == "" {
			return errors.New("remarks is required")
		}

		if r.Code <= 0 {
			return errors.New("code is required")
		}

	}

	return nil
}
