package requests

import "errors"

type EnrollmentRequest struct {
	StudentNumber     string `json:"student_number"`
	StudyPeriodNumber string `json:"study_period_number"`
	CampusNumber      string `json:"campus_number"`
	SchoolNumber      string `json:"school_number"`
	EnrollmentDate    string `json:"enrollment_date"`
	Status            string `json:"status"`
	CampusClassNumber string `json:"campus_class_number"`
	ContinuingStudent bool   `json:"continuing_student"`
}

func (e *EnrollmentRequest) Validate() error {
	if e.CampusNumber == "" {
		return errors.New("campus number is required")
	}
	if e.SchoolNumber == "" {
		return errors.New("school number is required")
	}
	if e.CampusClassNumber == "" {
		return errors.New("class is required")

	}
	if e.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if e.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}
	if e.EnrollmentDate == "" {
		return errors.New("enrollment date is required")
	}
	if e.Status == "" {
		return errors.New("status is required")
	}
	if e.Status != "active" && e.Status != "inactive" {
		return errors.New("status must be active or inactive")
	}

	return nil
}
