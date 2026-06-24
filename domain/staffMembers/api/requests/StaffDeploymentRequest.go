package requests

import (
	"errors"
)

type StaffDeploymentRequest struct {
	JobOfferNumber string `json:"job_offer_number"`
}

type EditStaffDeploymentRequest struct {
	Number string `json:"number"`
	Status string `json:"status"`
}

func (r *StaffDeploymentRequest) Validate() error {
	if r.JobOfferNumber == "" {
		return errors.New("job_offer_number is required")
	}
	return nil
}

func (r *EditStaffDeploymentRequest) Validate() error {
	if r.Number == "" {
		return errors.New("number is required")
	}
	if r.Status == "" {
		return errors.New("status is required")
	}
	return nil
}
