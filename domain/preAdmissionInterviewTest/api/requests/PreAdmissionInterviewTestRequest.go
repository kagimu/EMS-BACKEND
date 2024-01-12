package requests

import "errors"

type PreAdmissionInterviewTestRequest struct {
	StudentNumber   string `json:"student_number"`
	SubjectNumber   string `json:"subject_number"`
	Marks           int    `json:"marks"`
	GradeNumber     string `json:"grade_number"`
	ClassNumber     string `json:"class_number"`
	InterviewNumber string `json:"interview_number"`
}

func (p *PreAdmissionInterviewTestRequest) Validate() error {
	if p.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if p.SubjectNumber == "" {
		return errors.New("subject number is required")
	}
	if p.InterviewNumber == "" {
		return errors.New("location is required")
	}

	return nil
}
