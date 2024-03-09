package requests

import "errors"

type AssessmentAreaRequest struct {
	Name          string   `json:"name"`
	ClassNumber   []string `json:"class_number"`
	SubjectNumber string   `json:"subject_number"`
}

func (a *AssessmentAreaRequest) Validate() error {
	if len(a.ClassNumber) == 0 {
		return errors.New("at least one class number is required")
	}
	if a.SubjectNumber == "" {
		return errors.New("subject number is required")
	}
	if a.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
