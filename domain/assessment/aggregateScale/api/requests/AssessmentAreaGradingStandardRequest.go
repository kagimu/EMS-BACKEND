package requests

import "errors"

type AssessmentAreaGradingStandardRequest struct {
	AssessmentAreaNumber  string `json:"assessment_area_number"`
	GradingStandardNumber string `json:"grading_standard_number"`
	Code                  int64  `json:"code"`
}

func (asgs *AssessmentAreaGradingStandardRequest) Validate() error {

	if asgs.AssessmentAreaNumber == "" {
		return errors.New("assessment area  number is required")
	}

	if asgs.GradingStandardNumber == "" {
		return errors.New("grading standard number is required")
	}

	if asgs.Code == 0 {
		return errors.New("code is required")
	}

	return nil
}
