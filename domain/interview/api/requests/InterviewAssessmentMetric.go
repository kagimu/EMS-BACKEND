package requests

import "errors"

type InterviewAssessmentMetricRequest struct {
	InterviewNumber        string `json:"interview_number"`
	AssessmentMetricNumber string `json:"assessment_metric_number"`
}

func (a *InterviewAssessmentMetricRequest) Validate() error {
	if a.InterviewNumber == "" {
		return errors.New("interview number is required")
	}
	if a.AssessmentMetricNumber == "" {
		return errors.New("assessment metric number is required")
	}

	return nil
}
