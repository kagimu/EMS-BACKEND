package requests

import "errors"

type PromotionRequest struct {
	Condition                 string `json:"condition"`
	AssessmentNumber          string `json:"assessment_number"`
	ReasonForAltering         string `json:"reason_for_altering"`
	GradeNumber               string `json:"grade_number"`
	StudentNumber             string `json:"student_number"`
	AssessmentGradeNumber     string `json:"assessment_grade_number"`
	Authorizer                string `json:"authorizer"`
	Status                    string `json:"status"`
	Type                      string `json:"type"`
	PreviousCampusClassNumber string `json:"previous_campus_class_number"`
	NextCampusClassNumber     string `json:"next_campus_class_number"`
}

func (p *PromotionRequest) Validate() error {
	if p.AssessmentNumber == "" {
		return errors.New("assessment number is required")
	}
	if p.AssessmentGradeNumber == "" {
		return errors.New("assessment grade number is required")
	}
	if p.Authorizer == "" {
		return errors.New("authorizer is required")
	}
	if p.StudentNumber == "" {
		return errors.New("student number is required")
	}

	return nil
}
