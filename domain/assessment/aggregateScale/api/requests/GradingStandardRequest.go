package requests

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
)

type GradingStandardRequest struct {
	Code         int64                    `json:"code"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	SchoolNumber string                   `json:"school_number"`
	Status       string                   `json:"status"`
	Classes      []GradingStandardClasses `json:"classes"`
}

func (a *GradingStandardRequest) Validate() error {
	if a.Name == "" {
		return errors.New("name is required")
	}
	if a.SchoolNumber == "" {
		return errors.New("school number is required")
	}

	if a.Status != "" && a.Status != constants.StatusInactive && a.Status != constants.StatusActive {
		return errors.New("unknown status " + a.Status)
	}

	for _, class := range a.Classes {
		if class.ClassNumber == "" {
			return errors.New("choose a class")
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
	ClassNumber string `json:"class_number"`
	ScaleNumber string `json:"scale_number"`
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
