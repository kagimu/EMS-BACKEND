package requests

import "errors"

type AssessmentRequest struct {
	AssessmentMetricNumber string `json:"assessment_metric_number"`
	StudentNumber          string `json:"student_number"`
	AssessmentAreaNumber   string `json:"assessment_area_number"`
	LevelNumber            string `json:"level_number"`
	AssessmentDate         string `json:"assessment_date"`
}

func (a *AssessmentRequest) Validate() error {
	if a.AssessmentMetricNumber == "" {
		return errors.New("assessment metric number is required")
	}
	if a.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if a.AssessmentAreaNumber == "" {
		return errors.New("assessment area is required")
	}
	if a.LevelNumber == "" {
		return errors.New("level number is required")
	}
	if a.AssessmentDate == "" {
		return errors.New("assessment date is required")
	}

	return nil
}
