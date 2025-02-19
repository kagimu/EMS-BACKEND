package requests

import (
	"errors"
	"fmt"
	requests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/api/request"
)

type AggregateScaleRequest struct {
	AggregateScales      []AggregateScaleItem                  `json:"aggregate_scales"`
	ScaleName            string                                `json:"scale_name"`
	AggregateScaleCount  int                                   `json:"aggregate_scale_count"`
	SchoolNumber         string                                `json:"school_number"`
	AssessmentGrades     []requests.AssessmentGradeItemRequest `json:"assessment_grades"`
	AssessmentGradeCount int                                   `json:"assessment_grade_count"`
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

	if aggregateScale.AggregateScales == nil || len(aggregateScale.AggregateScales) == 0 {
		return errors.New("please set the aggregate scale")
	}

	if len(aggregateScale.AggregateScales) != aggregateScale.AggregateScaleCount {
		return errors.New("aggregate scales should be equal to the defined count")
	}

	if aggregateScale.SchoolNumber == "" {
		return errors.New("choose school")

	}

	if aggregateScale.ScaleName == "" {
		return errors.New("provide scale name")

	}

	var runsToEnd, startsAtZero bool

	scales := aggregateScale.AggregateScales

	for count, grade := range scales {

		var gradeIndex = count + 1

		if count == 0 {
			if grade.MinValue != 0 {
				return errors.New(fmt.Sprintf("the minimum value for the grade at count %d should be zero", gradeIndex))
			}
		}

		if gradeIndex == len(scales) {
			if grade.MaxValue != 100 {
				return errors.New(fmt.Sprintf("the maximum value for the grade at count %d should be 100", gradeIndex))
			}
		}

		if grade.Weight <= 0 {
			return errors.New(fmt.Sprintf("please set the weight for the grade at count %d", gradeIndex))
		}

		if grade.Remark == "" {
			return errors.New(fmt.Sprintf("please set the remark for the grade at count %d", gradeIndex))
		}

		if grade.GradeNumber == "" {
			return errors.New(fmt.Sprintf("please choose the assessment grade at count %d", gradeIndex))
		}

		if grade.Code <= 0 {
			return errors.New(fmt.Sprintf("please set the code for the grade at count %d", gradeIndex))
		}

		if grade.MinValue < 0 {
			return errors.New(fmt.Sprintf("minimum value must be greater than 0 or equal to zero for the grade at count %d", gradeIndex))
		}

		if grade.MinValue >= 100 {
			return errors.New(fmt.Sprintf("minimum value should be less than 100 for the grade at count %d", gradeIndex))
		}

		if grade.MaxValue <= 0 {
			return errors.New(fmt.Sprintf("maximum value must be greater than 0 for the grade at count %d", gradeIndex))
		}

		if grade.MaxValue > 100 {
			return errors.New(fmt.Sprintf("maximum value should not be greater than 100 for the grade at count %d", gradeIndex))
		}

		if grade.MaxValue <= grade.MinValue {
			return errors.New(fmt.Sprintf("maximum value must be greater than minimum value for the grade at count %d", gradeIndex))
		}

		if count > 0 {
			err := checkPrevRow(scales[count-1], grade, gradeIndex)
			if err != nil {
				return err
			}
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

	if aggregateScale.AssessmentGrades == nil || len(aggregateScale.AssessmentGrades) == 0 {
		return errors.New("set grading category for this scale")
	}

	if len(aggregateScale.AssessmentGrades) != aggregateScale.AssessmentGradeCount {
		return errors.New("grade categories should be equal to the defined count")
	}

	for _, r := range aggregateScale.AssessmentGrades {

		if r.Abbreviation == "" {
			return errors.New("abbreviation is required")
		}

		if r.Name == "" {
			return errors.New("name is required")
		}

		if r.Description == "" {
			return errors.New("description is required")
		}

		if r.MinimumPoints <= 0 {
			return errors.New("minimum points is required")
		}

		if r.MaximumPoints <= 0 {
			return errors.New("maximum points is required")
		}

		if r.Remarks == "" {
			return errors.New("remarks is required")
		}

		if r.Code <= 0 {
			return errors.New("code is required")
		}

	}

	return nil
}

func checkPrevRow(prevItem, currentItem AggregateScaleItem, currentIndex int) error {
	if currentItem.MinValue <= prevItem.MaxValue {
		return errors.New(fmt.Sprintf("The minimum score at grade count %d should be greator than the maximum score at grade count %d", currentIndex, currentIndex-1))
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
