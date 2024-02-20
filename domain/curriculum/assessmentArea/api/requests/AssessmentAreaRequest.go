package requests

import "errors"

type AssessmentAreaRequest struct {
	Name        string   `json:"name"`
	ClassNumber []string `json:"class_number"`
}

func (a *AssessmentAreaRequest) Validate() error {
	if len(a.ClassNumber) == 0 {
		return errors.New("class number is required")
	}
	if a.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
