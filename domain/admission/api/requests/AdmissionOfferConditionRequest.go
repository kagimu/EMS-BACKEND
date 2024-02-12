package requests

import "errors"

type AdmissionOfferConditionRequest struct {
	Description              string `json:"description"`
	CustomAdmissionCondition string `json:"custom_admission_condition"`
}

func (a *AdmissionOfferConditionRequest) Validate() error {
	if a.Description == "" {
		return errors.New("description is required")
	}

	return nil
}
