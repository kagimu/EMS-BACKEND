package requests

import "errors"

type JobPositionRequest struct {
	PostName         string `json:"post_name"`
	SchoolNumber     string `json:"school_number"`
	DepartmentNumber string `json:"department_number"`
}

type EditJobPositionRequest struct {
	Number           string `json:"number"`
	PostName         string `json:"post_name"`
	SchoolNumber     string `json:"school_number"`
	DepartmentNumber string `json:"department_number"`
}

func (r *JobPositionRequest) Validate() error {
	if r.PostName == "" {
		return errors.New("post_name is required")
	}
	if r.SchoolNumber == "" {
		return errors.New("school_number is required")
	}
	if r.DepartmentNumber == "" {
		return errors.New("department_number is required")
	}
	return nil
}

func (r *EditJobPositionRequest) Validate() error {
	if r.Number == "" {
		return errors.New("number is required")
	}
	if r.PostName == "" {
		return errors.New("post_name is required")
	}
	if r.SchoolNumber == "" {
		return errors.New("school_number is required")
	}
	if r.DepartmentNumber == "" {
		return errors.New("department_number is required")
	}
	return nil
}
