package requests

import "errors"

type AssessmentSubMetricRequest struct {
	Name                   string `json:"name"`
	AssessmentMetricNumber string `json:"assessment_metric_number"`
}

func (a *AssessmentSubMetricRequest) Validate() error {
	if a.Name == " " {
		return errors.New("name is required")
	}

	if a.AssessmentMetricNumber == "" {
		return errors.New("assessment metric number is required")
	}

	return nil
}
