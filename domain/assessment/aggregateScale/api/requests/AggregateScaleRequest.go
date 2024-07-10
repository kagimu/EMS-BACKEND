package requests

import "errors"

type AggregateScaleRequest struct {
	Grades       []AggregateScaleItem `json:"grades"`
	ScaleName    string               `json:"scale_name"`
	Count        int64                `json:"count"`
	SchoolNumber string               `json:"school_number"`
}

type AggregateScaleItem struct {
	Code         int64   `json:"code"`
	GradeNumber  string  `json:"grade_number"`
	MinValue     float64 `json:"min_value"`
	MaxValue     float64 `json:"max_value"`
	Weight       float64 `json:"weight"`
	Remark       string  `json:"remark"`
	Illustration string  `json:"illustration"`
}

type EditAggregateScaleRequest struct {
	AggregateScaleRequest
	ScaleNumber string `json:"scale_number"`
}

func (aggregateScale *AggregateScaleRequest) Validate() error {

	if aggregateScale.Grades == nil || len(aggregateScale.Grades) == 0 {
		return errors.New("please select the grades")
	}

	if aggregateScale.SchoolNumber == "" {
		return errors.New("choose school")

	}

	if aggregateScale.ScaleName == "" {
		return errors.New("provide scale name")

	}

	var runsToEnd, startsAtZero bool

	for _, grade := range aggregateScale.Grades {

		if grade.Weight <= 0 {
			return errors.New("set the weight for this score")
		}

		if grade.Remark == "" {
			return errors.New("set the remark for this score")
		}

		if grade.GradeNumber == "" {
			return errors.New("choose the assessment grade")
		}

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

		if grade.MaxValue >= 100 {
			runsToEnd = true
		}

		if grade.MinValue <= 0 {
			startsAtZero = true
		}

	}

	if !startsAtZero || !runsToEnd {
		return errors.New("the scale must start from 0 to 100")
	}

	return nil
}

func (aggregateScale *EditAggregateScaleRequest) Validate() error {
	if aggregateScale.ScaleNumber == "" {
		return errors.New("choose the scale to edit")

	}

	if err := aggregateScale.AggregateScaleRequest.Validate(); err != nil {
		return err
	}

	return nil
}
