package requests

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
)

type AssessmentMetricRequest struct {
	MaxScore               float64 `json:"max_score"`
	Description            string  `json:"description"`
	AssessmentAreaNumber   string  `json:"assessment_area_number"`
	StudyPeriodNumber      string  `json:"study_period_number"`
	ScheduledOn            string  `json:"scheduled_on"`
	MetricNumber           string  `json:"metric_number"`
	ClassNumber            string  `json:"class_number"`
	CampusNumber           string  `json:"campus_number"`
	CurriculumNumber       string  `json:"curriculum_number"`
	CampusCurriculumNumber string  `json:"campus_curriculum_number"`
	RequireComment         string  `json:"require_comment"`
}

type MultipleAssessmentScheduleRequest struct {
	Description            string                      `json:"description"`
	StudyPeriodNumber      string                      `json:"study_period_number"`
	MetricNumber           string                      `json:"metric_number"`
	ClassNumber            string                      `json:"class_number"`
	CampusNumber           string                      `json:"campus_number"`
	CurriculumNumber       string                      `json:"curriculum_number"`
	CampusCurriculumNumber string                      `json:"campus_curriculum_number"`
	Code                   int64                       `json:"code"`
	Schedules              []AssessmentScheduleRequest `json:"schedules"`
}

func (assessmentMetric MultipleAssessmentScheduleRequest) Validate() error {

	if assessmentMetric.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}
	if assessmentMetric.MetricNumber == "" {
		return errors.New("metric number is required")
	}
	if assessmentMetric.CampusNumber == "" {
		return errors.New("Choose campus ")
	}
	if assessmentMetric.ClassNumber == "" {
		return errors.New("Choose a class ")
	}

	schedules := assessmentMetric.Schedules
	if len(schedules) == 0 {
		return errors.New("No assessment schedules, set at least one schedule ")
	}

	for _, schedule := range schedules {

		if schedule.AssessmentAreaNumber == "" {
			return errors.New("assessment area number is required")
		}

		if schedule.ScheduledOn == "" {
			return errors.New("scheduled on is required")
		}
		if schedule.RequireComment == "" {
			return errors.New("please specify if the comment  is required or not")
		}
		if schedule.RequireComment != constants.StatusYes && schedule.RequireComment != constants.StatusNo {
			return errors.New("require comment should be YES or NO")
		}
		if schedule.MaxScore <= 0 {
			return errors.New("max score can not be zero or negative")
		}
	}
	
	return checkAssessmentAreaDuplicatesInAssessmentScheduleRequest(schedules)
}

func checkAssessmentAreaDuplicatesInAssessmentScheduleRequest(assessmentScheduleRequests []AssessmentScheduleRequest) error {

	cleanSelections := make([]AssessmentScheduleRequest, len(assessmentScheduleRequests))

	for _, assessmentScheduleRequest := range assessmentScheduleRequests {

		for _, cleanSelection := range cleanSelections {

			if cleanSelection.AssessmentAreaNumber == assessmentScheduleRequest.AssessmentAreaNumber {

				assessmentAreaByNumber, err := repositories.GetAssessmentAreaByNumber(assessmentScheduleRequest.AssessmentAreaNumber)
				if err != nil {
					return errors.New("error fetching assessment area while checking duplicates: " + err.Error())
				}

				return errors.New(fmt.Sprintf("%s is selected multiple times", assessmentAreaByNumber.Name))
			}

		}

		cleanSelections = append(cleanSelections, assessmentScheduleRequest)
	}

	return nil
}

type AssessmentScheduleRequest struct {
	MaxScore             float64 `json:"max_score"`
	AssessmentAreaNumber string  `json:"assessment_area_number"`
	RequireComment       string  `json:"require_comment"`
	ScheduledOn          string  `json:"scheduled_on"`
}

type UpdateAssessmentMetricRequest struct {
	AssessmentMetricRequest
	AssessmentMetricNumber string `json:"assessment_metric_number"`
}

func (assessmentMetric *AssessmentMetricRequest) Validate() error {
	if assessmentMetric.AssessmentAreaNumber == "" {
		return errors.New("assessment area number is required")
	}
	if assessmentMetric.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}
	if assessmentMetric.MetricNumber == "" {
		return errors.New("metric number is required")
	}
	if assessmentMetric.CampusNumber == "" {
		return errors.New("Choose campus ")
	}
	if assessmentMetric.ClassNumber == "" {
		return errors.New("Choose a class ")
	}
	if assessmentMetric.ScheduledOn == "" {
		return errors.New("scheduled on is required")
	}
	if assessmentMetric.RequireComment == "" {
		return errors.New("please specify if the comment  is required or not")
	}
	if assessmentMetric.RequireComment != constants.StatusYes && assessmentMetric.RequireComment != constants.StatusNo {
		return errors.New("require comment should be YES or NO")
	}
	if assessmentMetric.MaxScore <= 0 {
		return errors.New("max score can not be zero or negative")
	}

	return nil
}

func (assessmentMetric *UpdateAssessmentMetricRequest) Validate() error {

	if assessmentMetric.AssessmentMetricNumber == "" {
		return errors.New("choose the assessment metric being edited")
	}

	err := assessmentMetric.AssessmentMetricRequest.Validate()
	if err != nil {
		return err
	}

	return nil
}
