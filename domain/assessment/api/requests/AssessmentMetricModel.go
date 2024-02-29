package requests

import "errors"

type AssessmentMetricRequest struct {
	Name                 string `json:"name"`
	Weight               uint8  `json:"weight"`
	DescriptiveScore     string `json:"descriptive_score"`
	PercentageScore      uint8  `json:"percentage_score"`
	MaxScore             uint8  `json:"max_score"`
	ActualScore          uint8  `json:"actual_score"`
	Description          string `json:"description"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
}

func (a *AssessmentMetricRequest) Validate() error {
	if a.Name == "" {
		return errors.New("name is required")
	}
	if a.Weight == 0 {
		return errors.New("weight is required")
	}

	return nil
}
