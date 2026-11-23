package requests

import (
	"errors"
	"time"
)

type JobVacancyRequest struct {
	JobPositionNumber string `json:"job_position_number"`
	VacancyTitle      string `json:"vacancy_title"`
	ReferenceNumber   string `json:"reference_number"`
	DepartmentNumber  string `json:"department_number"`
	SchoolNumber      string `json:"school_number"`
	BranchLocation    string `json:"branch_location"`
	NumberOfPositions int    `json:"number_of_positions"`
	AppointmentTerms  string `json:"appointment_terms"`
	OpeningDate       string `json:"opening_date"`
	ClosingDate       string `json:"closing_date"`
	Status            string `json:"status"`
}

type EditJobVacancyRequest struct {
	Number            string `json:"number"`
	JobPositionNumber string `json:"job_position_number"`
	VacancyTitle      string `json:"vacancy_title"`
	ReferenceNumber   string `json:"reference_number"`
	DepartmentNumber  string `json:"department_number"`
	SchoolNumber      string `json:"school_number"`
	BranchLocation    string `json:"branch_location"`
	NumberOfPositions int    `json:"number_of_positions"`
	AppointmentTerms  string `json:"appointment_terms"`
	OpeningDate       string `json:"opening_date"`
	ClosingDate       string `json:"closing_date"`
	Status            string `json:"status"`
}

func (r *JobVacancyRequest) Validate() error {
	if r.JobPositionNumber == "" {
		return errors.New("job_position_number is required")
	}
	if r.VacancyTitle == "" {
		return errors.New("vacancy_title is required")
	}
	if r.SchoolNumber == "" {
		return errors.New("school_number is required")
	}
	if r.OpeningDate != "" {
		if _, err := time.Parse("2006-01-02", r.OpeningDate); err != nil {
			return errors.New("invalid opening_date format, must be YYYY-MM-DD")
		}
	}
	if r.ClosingDate != "" {
		if _, err := time.Parse("2006-01-02", r.ClosingDate); err != nil {
			return errors.New("invalid closing_date format, must be YYYY-MM-DD")
		}
	}
	return nil
}

func (r *EditJobVacancyRequest) Validate() error {
	if r.Number == "" {
		return errors.New("number is required")
	}
	if r.JobPositionNumber == "" {
		return errors.New("job_position_number is required")
	}
	if r.VacancyTitle == "" {
		return errors.New("vacancy_title is required")
	}
	if r.SchoolNumber == "" {
		return errors.New("school_number is required")
	}
	if r.OpeningDate != "" {
		if _, err := time.Parse("2006-01-02", r.OpeningDate); err != nil {
			return errors.New("invalid opening_date format, must be YYYY-MM-DD")
		}
	}
	if r.ClosingDate != "" {
		if _, err := time.Parse("2006-01-02", r.ClosingDate); err != nil {
			return errors.New("invalid closing_date format, must be YYYY-MM-DD")
		}
	}
	return nil
}
