package requests

import (
	"errors"
	"html"
	"strings"
)

type AdmissionRequest struct {
	Status            string `json:"status"`
	StudentNumber     string `json:"student_number"`
	AdmissionFeePaid  bool   `json:"admission_fee_paid"`
	InterviewNumber   string `json:"interview_number"`
	StudyPeriodNumber string `json:"study_period_number"`
	ClassNumber       string `json:"class_number"`
	StreamNumber      string `json:"stream_number"`
	CampusNumber      string `json:"campus_number"`
}

func (a *AdmissionRequest) Validate() error {
	if a.Status == "" {
		return errors.New("status is required")
	}
	if a.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if a.CampusNumber == "" {
		return errors.New("campus number is required")
	}

	return nil
}

func (a *AdmissionRequest) Prepare() {
	a.StudentNumber = html.EscapeString(strings.TrimSpace(a.StudentNumber))
	a.ClassNumber = html.EscapeString(strings.TrimSpace(a.ClassNumber))
	a.StreamNumber = html.EscapeString(strings.TrimSpace(a.StreamNumber))
	a.CampusNumber = html.EscapeString(strings.TrimSpace(a.CampusNumber))
}
