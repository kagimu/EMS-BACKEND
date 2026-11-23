package requests

import (
	"errors"
)

type JobOfferRequest struct {
	ApplicationNumber     string `json:"application_number"`
	InterviewNumber       string `json:"interview_number"`
	AppointmentLetterHtml string `json:"appointment_letter_html"`
	Status                string `json:"status"`
}

func (r *JobOfferRequest) Validate() error {
	if r.ApplicationNumber == "" {
		return errors.New("application_number is required")
	}
	if r.InterviewNumber == "" {
		return errors.New("interview_number is required")
	}
	if r.AppointmentLetterHtml == "" {
		return errors.New("appointment_letter_html is required")
	}
	if r.Status == "" {
		return errors.New("status is required")
	}
	return nil
}
