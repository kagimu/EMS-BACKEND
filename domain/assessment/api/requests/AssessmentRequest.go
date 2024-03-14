package requests

import "errors"

type AssessmentRequest struct {
	AssessmentMetricNumber string  `json:"assessment_metric_number"`
	StudentNumber          string  `json:"student_number"`
	AssessmentAreaNumber   string  `json:"assessment_area_number"`
	LevelNumber            string  `json:"level_number"`
	Comment                string  `json:"comment"`
	AssessmentDate         string  `json:"assessment_date"`
	ActualScore            float64 `json:"actual_score"`
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
