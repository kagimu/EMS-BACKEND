package requests

import "errors"

type AdmissionOfferRequest struct {
	ClassNumber  string `json:"class_number"`
	SchoolNumber string `json:"school_number"`
	CampusNumber string `json:"campus_number"`
	StudentAPN   string `json:"student_apn"`
}

func (a *AdmissionOfferRequest) Validate() error {
	if a.CampusNumber == "" {
		return errors.New("campus number is required")
	}
	if a.ClassNumber == "" {
		return errors.New("class number is required")
	}
	if a.SchoolNumber == "" {
		return errors.New("school number is required")
	}
	if a.StudentAPN == "" {
		return errors.New("student apn is required")
	}

	return nil
}
