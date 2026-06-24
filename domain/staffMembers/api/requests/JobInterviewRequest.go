package requests

import (
	"errors"
)

type JobInterviewRequest struct {
	JobVacancyNumber string `json:"job_vacancy_number"`
	InterviewTitle   string `json:"interview_title"`
	InterviewDate    string `json:"interview_date"`
	InterviewTime    string `json:"interview_time"`
	LocationOrLink   string `json:"location_or_link"`
	Status           string `json:"status"`
}

type EditJobInterviewRequest struct {
	Number           string `json:"number"`
	JobVacancyNumber string `json:"job_vacancy_number"`
	InterviewTitle   string `json:"interview_title"`
	InterviewDate    string `json:"interview_date"`
	InterviewTime    string `json:"interview_time"`
	LocationOrLink   string `json:"location_or_link"`
	Status           string `json:"status"`
}

type ScheduleInterviewRequest struct {
	ApplicationNumber string `json:"number"`
	InterviewNumber   string `json:"interview_number"`
	Remarks           string `json:"remarks"`
}

func (r *JobInterviewRequest) Validate() error {
	if r.JobVacancyNumber == "" {
		return errors.New("job_vacancy_number is required")
	}
	if r.InterviewTitle == "" {
		return errors.New("interview_title is required")
	}
	if r.InterviewDate == "" {
		return errors.New("interview_date is required")
	}
	return nil
}

func (r *EditJobInterviewRequest) Validate() error {
	if r.Number == "" {
		return errors.New("number is required")
	}
	if r.JobVacancyNumber == "" {
		return errors.New("job_vacancy_number is required")
	}
	if r.InterviewTitle == "" {
		return errors.New("interview_title is required")
	}
	if r.InterviewDate == "" {
		return errors.New("interview_date is required")
	}
	return nil
}

func (r *ScheduleInterviewRequest) Validate() error {
	if r.ApplicationNumber == "" {
		return errors.New("application number is required")
	}
	if r.InterviewNumber == "" {
		return errors.New("interview_number is required")
	}
	return nil
}

type ScoreInterviewRequest struct {
	ScheduledInterviewNumber string `json:"scheduled_interview_number"`
	Comment                  string `json:"comment"`
	Status                   string `json:"status"`
}

func (r *ScoreInterviewRequest) Validate() error {
	if r.ScheduledInterviewNumber == "" {
		return errors.New("scheduled_interview_number is required")
	}
	if r.Status == "" {
		return errors.New("status is required")
	}
	return nil
}
