package requests

import "errors"

type DepartmentRequest struct {
	Name         string `json:"name"`
	Code         int64  `json:"code"`
	SchoolNumber string `json:"school_number"`
}

type EditDepartmentRequest struct {
	Number       string `json:"number"`
	Name         string `json:"name"`
	Code         int64  `json:"code"`
	SchoolNumber string `json:"school_number"`
}

func (d *DepartmentRequest) Validate() error {
	if d.Name == "" {
		return errors.New("name is required")
	}
	if d.SchoolNumber == "" {
		return errors.New("school_number is required")
	}
	return nil
}

func (d *EditDepartmentRequest) Validate() error {
	if d.Number == "" {
		return errors.New("number is required")
	}
	if d.Name == "" {
		return errors.New("name is required")
	}
	if d.SchoolNumber == "" {
		return errors.New("school_number is required")
	}
	return nil
}
