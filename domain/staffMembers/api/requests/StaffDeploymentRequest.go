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

type StaffRoleItem struct {
	RoleNumber string `json:"role_number"`
	ExpiryDate string `json:"expiry_date"`
}

type AssignStaffRoleRequest struct {
	DeploymentNumber string          `json:"deployment_number"`
	Roles            []StaffRoleItem `json:"roles"`
	UserStatus       string          `json:"user_status"`
}

func (r *AssignStaffRoleRequest) Validate() error {
	if r.DeploymentNumber == "" {
		return errors.New("deployment_number is required")
	}
	for _, role := range r.Roles {
		if role.RoleNumber == "" {
			return errors.New("role_number is required for each role")
		}
		if role.ExpiryDate == "" {
			return errors.New("expiry_date is required for each role")
		}
	}
	return nil
}

