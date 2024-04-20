package requests

import "errors"

type GradingStandardRequest struct {
	Code              int64  `json:"code"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Type              string `json:"type"`
	CampusClassNumber string `json:"campus_class_number"`
}

func (a *GradingStandardRequest) Validate() error {
	if a.Name == "" {
		return errors.New("name is required")
	}
	if a.Type == "" {
		return errors.New("type is required")
	}
	if a.CampusClassNumber == "" {
		return errors.New("campus class number is required")
	}
	return nil
}
