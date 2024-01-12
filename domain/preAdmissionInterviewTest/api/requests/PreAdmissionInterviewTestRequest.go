package requests

import "errors"

type PreAdmissionInterviewTestRequest struct {
	StudentNumber string `json:"student_number"`
	SubjectNumber string `json:"subject_number"`
	Marks         int    `json:"marks"`
	GradeNumber   string `json:"grade_number"`
	ClassNumber   string `json:"class_number"`
	Location      string `json:"location"`
	InterviewDate string `json:"interview_date"`
}

func (p *PreAdmissionInterviewTestRequest) Validate() error {
	if p.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if p.SubjectNumber == "" {
		return errors.New("subject number is required")
	}
	if p.Location == "" {
		return errors.New("location is required")
	}
	if p.InterviewDate == "" {
		return errors.New("interview date is required")
	}

	return nil
}
