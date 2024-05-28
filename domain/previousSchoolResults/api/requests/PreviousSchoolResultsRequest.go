package requests

import "errors"

type PreviousSchoolResultsRequest struct {
	StudentNumber      string               `json:"student_number"`
	ApplicationNumber  string               `json:"application_number"`
	AssessmentAndScore []AssessmentAndScore `json:"assessment_and_score"`
}

type AssessmentAndScore struct {
	AssessmentAreaNumber string  `json:"assessment_area_number"`
	Score                float64 `json:"score"`
}

func (p *PreviousSchoolResultsRequest) Validate() error {
	if p.StudentNumber == "" {
		return errors.New("student number is required")
	}

	return nil
}
