package requests

import "errors"

type AssessmentRequest struct {
	AssessmentMetricNumber string  `json:"assessment_metric_number"`
	StudentNumber          string  `json:"student_number"`
	Comment                string  `json:"comment"`
	ActualScore            float64 `json:"actual_score"`
}

func (a *AssessmentRequest) Validate() error {
	if a.AssessmentMetricNumber == "" {
		return errors.New("assessment metric number is required")
	}
	if a.StudentNumber == "" {
		return errors.New("student number is required")
	}
	return nil
}
