package requests

import "errors"

type AdmissionOfferRequest struct {
	ClassNumber                   string   `json:"class_number"`
	SchoolNumber                  string   `json:"school_number"`
	CampusNumber                  string   `json:"campus_number"`
	StudentNumber                 string   `json:"student_number"`
	AdmissionOfferConditionNumber []string `json:"admission_offer_condition_number"`
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
	if a.StudentNumber == "" {
		return errors.New("student number is required")
	}

	return nil
}
