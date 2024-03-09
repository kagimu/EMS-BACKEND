package requests

import "errors"

type PromotionRequest struct {
	Condition             string `json:"condition"`
	AssessmentNumber      string `json:"assessment_number"`
	ReasonForAltering     string `json:"reason_for_altering"`
	GradeNumber           string `json:"grade_number"`
	StudentNumber         string `json:"student_number"`
	AssessmentGradeNumber string `json:"assessment_grade_number"`
	AuthorizedByNumber    string `json:"authorized_by_number"`
}

func (p *PromotionRequest) Validate() error {
	if p.AssessmentNumber == "" {
		return errors.New("assessment number is required")
	}
	if p.AssessmentGradeNumber == "" {
		return errors.New("assessment grade number is required")
	}
	if p.AuthorizedByNumber == "" {
		return errors.New("staff number is required")
	}

	return nil
}
