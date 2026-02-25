package requests

import (
	"errors"
)

type AdmissionLetterRequest struct {
	SchoolNumber      string `json:"school_number"`
	StudentNumber     string `json:"student_number"`
	CampusNumber      string `json:"campus_number"`
	CampusClassNumber string `json:"campus_class_number"`
	Phrase            string `json:"phrase"`
	PrincipalSignature string `json:"principal_signature"`
}

func (r *AdmissionLetterRequest) Validate() error {
	if r.SchoolNumber == "" {
		return errors.New("school number is required")
	}
	if r.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if r.CampusNumber == "" {
		return errors.New("campus number is required")
	}
	if r.CampusClassNumber == "" {
		return errors.New("campus class number is required")
	}
	return nil
}

type UpdateAdmissionLetterRequest struct {
	Number string `json:"number"`
	AdmissionLetterRequest
}

func (r *UpdateAdmissionLetterRequest) Validate() error {
	if r.Number == "" {
		return errors.New("admission letter number is required for update")
	}
	return r.AdmissionLetterRequest.Validate()
}
