package requests

import (
	"errors"
)

type NewTemplateRequest struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	SchoolNumber string `json:"school_number"`
}

func (newTemplateRequest *NewTemplateRequest) Validate() error {

	if newTemplateRequest.SchoolNumber == "" {
		return errors.New("choose school where to add this template")
	}

	if newTemplateRequest.Name == "" {
		return errors.New("template name is required")
	}

	if newTemplateRequest.Status == "" {
		return errors.New("choose the template status")
	}

	if len(newTemplateRequest.Name) > 250 {
		return errors.New("the name character length should not exceed 250")
	}

	if newTemplateRequest.Status != "ACTIVE" && newTemplateRequest.Status != "INACTIVE" {
		return errors.New("invalid template status, should be ACTIVE or INACTIVE")
	}

	return nil
}
