package requests

import (
	"errors"
)

type InterviewRequest struct {
	Message           string   `json:"message"`
	Venue             string   `json:"venue"`
	InterviewDate     string   `json:"interview_date"`
	Name              string   `json:"name"`
	CampusClassNumber string   `json:"campus_class_number"`
	AssessmentMetrics []string `json:"assessment_metrics"`
}

type EditInterviewRequest struct {
	InterviewRequest
	InterviewNumber string `json:"interview_number"`
}

func (i *InterviewRequest) Validate() error {
	if i.Message == "" {
		return errors.New("message is required")
	}
	if i.Venue == "" {
		return errors.New("venue is required")
	}
	if i.InterviewDate == "" {
		return errors.New("interview date is required")
	}
	if i.Name == "" {
		return errors.New("interview date is required")
	}

	return nil
}

func (i *EditInterviewRequest) Validate() error {

	if i.InterviewNumber == "" {
		return errors.New("choose interview being edited")
	}

	if err := i.InterviewRequest.Validate(); err != nil {
		return err
	}

	return nil
}
