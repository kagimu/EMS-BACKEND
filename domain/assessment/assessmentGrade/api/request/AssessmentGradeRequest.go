package requests

import "errors"

type AssessmentGradeRequest struct {
	SchoolNumber string `json:"school_number"`
	Code         int64  `json:"code"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

type EditAssessmentGradeRequest struct {
	AssessmentGradeRequest
	AssessmentGradeNumber string `json:"assessment_grade_number"`
}

func (r *AssessmentGradeRequest) Validate() error {
	if r.Abbreviation == "" {
		return errors.New("abbreviation is required")
	}
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Code == 0 {
		return errors.New("code is required")
	}
	return nil
}
