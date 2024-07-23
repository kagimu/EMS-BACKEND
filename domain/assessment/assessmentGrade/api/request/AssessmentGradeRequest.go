package requests

import (
	"errors"
)

type AssessmentGradeRequest struct {
	Classes              []string                     `json:"classes"`
	AssessmentAreaCount  int                          `json:"assessment_area_count"`
	AssessmentGradeItems []AssessmentGradeItemRequest `json:"assessment_grade_items"`
	SchoolNumber         string                       `json:"school_number"`
}

type EditAssessmentGradeRequest struct {
	AssessmentGradeRequest
	Status                string `json:"status"`
	AssessmentGradeNumber string `json:"assessment_grade_number"`
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

	if assessmentRequest.AssessmentAreaCount <= 0 {
		return errors.New("set the number of assessment areas to consider I.E best 8")
	}

	if assessmentRequest.SchoolNumber == "" {
		return errors.New("choose a school")
	}

	if assessmentRequest.AssessmentAreaCount != len(assessmentRequest.AssessmentGradeItems) {
		return errors.New("the scale must be equal to the count")
	}

	if len(assessmentRequest.Classes) <= 0 {
		return errors.New("choose classes")
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

func (assessmentRequest *EditAssessmentGradeRequest) Validate() error {

	if assessmentRequest.AssessmentGradeNumber == "" {
		return errors.New("set assessment grade number")
	}

	if err := assessmentRequest.AssessmentGradeRequest.Validate(); err != nil {
		return err
	}

	return nil
}
