package requests

import "errors"

type AssessmentAreaGradingStandardRequest struct {
	AssessmentAreaNumber  string `json:"assessment_area_number"`
	GradingStandardNumber string `json:"grading_standard_number"`
}

func (asgs *AssessmentAreaGradingStandardRequest) Validate() error {

	if asgs.AssessmentAreaNumber == "" {
		return errors.New("assessment area  number is required")
	}

	if asgs.GradingStandardNumber == "" {
		return errors.New("grading standard number is required")
	}

	return nil
}
