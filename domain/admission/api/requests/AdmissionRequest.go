package requests

import (
	"errors"
	"html"
	"strings"
)

type AdmissionRequest struct {
	Status           string `json:"status"`
	StudentNumber    string `json:"student_number"`
	AdmissionFeePaid bool   `json:"admission_fee_paid"`
	ClassNumber      string `json:"class_number"`
	StreamNumber     string `json:"stream_number"`
	CampusNumber     string `json:"campus_number"`
}

func (a *AdmissionRequest) Validate() error {
	if a.Status == "" {
		return errors.New("status is required")
	}
	if a.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if a.ClassNumber == "" {
		return errors.New("class number is required")
	}
	if a.StreamNumber == "" {
		return errors.New("stream number is required")
	}
	if a.CampusNumber == "" {
		return errors.New("campus number is required")
	}

	return nil
}

func (a *AdmissionRequest) Prepare() {
	a.StudentNumber = html.EscapeString(strings.TrimSpace(strings.ToLower(a.StudentNumber)))
	a.ClassNumber = html.EscapeString(strings.TrimSpace(strings.ToLower(a.ClassNumber)))
	a.StreamNumber = html.EscapeString(strings.TrimSpace(strings.ToLower(a.StreamNumber)))
	a.CampusNumber = html.EscapeString(strings.TrimSpace(strings.ToLower(a.CampusNumber)))
}
