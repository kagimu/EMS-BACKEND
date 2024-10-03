package requests

import "errors"

type AssessmentAreaGradingStandardRequest struct {
	AssessmentAreaNumber  string `json:"assessment_area_number"`
	SchoolNumber          string `json:"school_number"`
	GradingStandardNumber string `json:"grading_standard_number"`
	Code                  int64  `json:"code"`
}

func (assessmentAreaGradingStandard *AssessmentAreaGradingStandardRequest) Validate() error {

	if assessmentAreaGradingStandard.SchoolNumber == "" {
		return errors.New("choose the school")
	}

	if assessmentAreaGradingStandard.AssessmentAreaNumber == "" {
		return errors.New("assessment area  number is required")
	}

	if assessmentAreaGradingStandard.GradingStandardNumber == "" {
		return errors.New("grading standard number is required")
	}

	if assessmentAreaGradingStandard.Code == 0 {
		return errors.New("code is required")
	}

	return nil
}
