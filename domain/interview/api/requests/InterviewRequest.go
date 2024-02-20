package requests

import (
	"errors"
	"html"
	"strings"
)

type InterviewRequest struct {
	Message       string `json:"message"`
	Venue         string `json:"venue"`
	InterviewDate string `json:"interview_date"`
	Name          string `json:"name"`
	StudentNumber string `json:"student_number"`
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

func (i *InterviewRequest) Prepare() {
	i.Message = html.EscapeString(strings.TrimSpace((i.Message)))
	i.Venue = html.EscapeString(strings.TrimSpace((i.Venue)))
}
