package requests

import (
	"errors"
)

type JobApplicationRequest struct {
	JobVacancyNumber string `json:"job_vacancy_number"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
}

type EditJobApplicationRequest struct {
	Number           string `json:"number"`
	JobVacancyNumber string `json:"job_vacancy_number"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	Status           string `json:"status"`
}

func (r *JobApplicationRequest) Validate() error {
	if r.JobVacancyNumber == "" {
		return errors.New("job_vacancy_number is required")
	}
	if r.FirstName == "" {
		return errors.New("first_name is required")
	}
	if r.LastName == "" {
		return errors.New("last_name is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.PhoneNumber == "" {
		return errors.New("phone_number is required")
	}
	return nil
}

func (r *EditJobApplicationRequest) Validate() error {
	if r.Number == "" {
		return errors.New("number is required")
	}
	if r.JobVacancyNumber == "" {
		return errors.New("job_vacancy_number is required")
	}
	if r.FirstName == "" {
		return errors.New("first_name is required")
	}
	if r.LastName == "" {
		return errors.New("last_name is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.PhoneNumber == "" {
		return errors.New("phone_number is required")
	}
	if r.Status == "" {
		return errors.New("status is required")
	}
	return nil
}
