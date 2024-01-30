package requests

import (
	"errors"
)

type InterviewStudentRequest struct {
	StudentNumber   string `json:"student_number"`
	InterviewNumber string `json:"interview_number"`
}

func (i *InterviewStudentRequest) Validate() error {
	if i.StudentNumber == "" {
		return errors.New("student is required")
	}
	if i.InterviewNumber == "" {
		return errors.New("venue is required")
	}

	return nil
}
