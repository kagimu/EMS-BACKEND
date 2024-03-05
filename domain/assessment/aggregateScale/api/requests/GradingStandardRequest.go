package requests

import "errors"

type GradingStandardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	ClassNumber string `json:"class_number"`
}

func (a *GradingStandardRequest) Validate() error {
	if a.Name == "" {
		return errors.New("name is required")
	}
	if a.Type == "" {
		return errors.New("type is required")
	}
	if a.ClassNumber == "" {
		return errors.New("class number is required")
	}
	return nil
}
