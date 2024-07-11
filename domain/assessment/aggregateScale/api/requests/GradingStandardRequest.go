package requests

import "errors"

type GradingStandardRequest struct {
	Code         int64                    `json:"code"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	SchoolNumber string                   `json:"school_number"`
	Classes      []GradingStandardClasses `json:"classes"`
}

func (a *GradingStandardRequest) Validate() error {
	if a.Name == "" {
		return errors.New("name is required")
	}
	if a.SchoolNumber == "" {
		return errors.New("school number is required")
	}

	for _, class := range a.Classes {
		if class.CampusClassNumber == "" {
			return errors.New("campus class number is required")
		}
		if class.ScaleNumber == "" {
			return errors.New("choose scale for the selected class")
		}
	}

	return nil
}

type UpdateGradingStandardRequest struct {
	GradingStandardRequest
	GradingStandardNumber string `json:"grading_standard_number"`
}

type GradingStandardClasses struct {
	CampusClassNumber string `json:"campus_class_number"`
	ScaleNumber       string `json:"scale_number"`
}

func (a *UpdateGradingStandardRequest) Validate() error {

	if err := a.GradingStandardRequest.Validate(); err != nil {
		return err
	}

	if a.GradingStandardNumber == "" {
		return errors.New("provide grading standard number to edit")
	}

	return nil
}
