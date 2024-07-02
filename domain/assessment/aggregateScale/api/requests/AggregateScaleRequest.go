package requests

import "errors"

type AggregateScaleRequest struct {
	Grades                []AggregateScaleItem `json:"grades"`
	GradingStandardNumber string               `json:"grading_standard_number"`
}

type AggregateScaleItem struct {
	Code        int64   `json:"code"`
	GradeNumber string  `json:"grade_number"`
	MinValue    float64 `json:"min_value"`
	MaxValue    float64 `json:"max_value"`
}

type EditAggregateScaleRequest struct {
	Code                  int64   `json:"code"`
	MinValue              float64 `json:"min_value"`
	MaxValue              float64 `json:"max_value"`
	CampusNumber          string  `json:"campus_number"`
	Illustration          string  `json:"illustration"`
	GradingStandardNumber string  `json:"grading_standard_number"`
	CampusClassNumber     string  `json:"campus_class_number"`
}

func (aggregateScale *AggregateScaleRequest) Validate() error {

	if aggregateScale.Grades == nil || len(aggregateScale.Grades) == 0 {
		return errors.New("please select the grades")
	}

	if aggregateScale.GradingStandardNumber == "" {
		return errors.New("grading standard number is required")

	}

	for _, grade := range aggregateScale.Grades {

		if grade.Code <= 0 {
			return errors.New("code is required")
		}

		if grade.MinValue < 0 {
			return errors.New("minimum value must be greater than 0 or equal to zero")
		}

		if grade.MinValue >= 100 {
			return errors.New("minimum value should be less than 100")
		}

		if grade.MaxValue <= 0 {
			return errors.New("maximum value must be greater than 0")
		}

		if grade.MaxValue > 100 {
			return errors.New("maximum value should not be greater than 100")
		}

		if grade.MaxValue < grade.MinValue {
			return errors.New("maximum value must be greater than minimum value")
		}

	}

	return nil
}
