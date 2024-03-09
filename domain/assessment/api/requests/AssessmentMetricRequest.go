package requests

import (
	"errors"
)

type AssessmentMetricRequest struct {
	Name                  string `json:"name"`
	Weight                uint8  `json:"weight"`
	MaxScore              uint8  `json:"max_score"`
	ActualScore           uint8  `json:"actual_score"`
	GradingStandardNumber string `json:"grading_standard_number"`
	Description           string `json:"description"`
	AssessmentAreaNumber  string `json:"assessment_area_number"`
	StudyPeriodNumber     string `json:"study_period_number"`
	CurriculumClassNumber string `json:"curriculum_class_number"`
}

func (a *AssessmentMetricRequest) Validate() error {
	if a.Name == "" {
		return errors.New("name is required")
	}
	if a.Weight == 0 {
		return errors.New("weight is required")
	}
	if a.GradingStandardNumber == "" {
		return errors.New("grading standard number is required")
	}
	if a.AssessmentAreaNumber == "" {
		return errors.New("assessment area number is required")
	}
	if a.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}
	if a.CurriculumClassNumber == "" {
		return errors.New("curriculum class number is required")
	}

	return nil
}
