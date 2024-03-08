package requests

import "errors"

type SalaryScaleRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	MinimumAmount uint64 `json:"minimum_amount"`
	MaximumAmount uint64 `json:"maximum_amount"`
}

type UpdateSalaryScaleRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	MinimumAmount uint64 `json:"minimum_amount"`
	MaximumAmount uint64 `json:"maximum_amount"`
	Number        string `json:"number"`
}

func (r *SalaryScaleRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.MinimumAmount == 0 {
		return errors.New("minimum amount is required")
	}
	if r.MaximumAmount == 0 {
		return errors.New("maximum amount is required")
	}
	if r.MaximumAmount < r.MinimumAmount {
		return errors.New("maximum amount must be greater than minimum amount")
	}

	return nil
}

func (r *UpdateSalaryScaleRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.MinimumAmount == 0 {
		return errors.New("minimum amount is required")
	}
	if r.MaximumAmount == 0 {
		return errors.New("maximum amount is required")
	}
	if r.MaximumAmount < r.MinimumAmount {
		return errors.New("maximum amount must be greater than minimum amount")
	}

	return nil
}
