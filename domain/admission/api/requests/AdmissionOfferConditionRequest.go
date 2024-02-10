package requests

import "errors"

type AdmissionOfferConditionRequest struct {
	Description string `json:"description"`
}

func (a *AdmissionOfferConditionRequest) Validate() error {
	if a.Description == "" {
		return errors.New("description is required")
	}

	return nil
}
