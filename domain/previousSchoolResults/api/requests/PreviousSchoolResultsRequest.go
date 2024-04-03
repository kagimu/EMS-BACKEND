package requests

import "errors"

type PreviousSchoolResultsRequest struct {
	AssessmentAreaNumber []string  `json:"assessment_area_number"`
	Score                []float64 `json:"score"`
	StudentNumber        string    `json:"student_number"`
	ApplicationNumber    string    `json:"application_number"`
}

func (p *PreviousSchoolResultsRequest) Validate() error {
	if p.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if len(p.AssessmentAreaNumber) == 0 {
		return errors.New("at least one assessment area number is required")
	}
	if len(p.Score) == 0 {
		return errors.New("at least one score is required")
	}
	return nil
}
