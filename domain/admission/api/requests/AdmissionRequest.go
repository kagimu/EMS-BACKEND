package requests

import (
	"errors"
	"html"
	"strings"
)

type AdmissionRequest struct {
	StudentNumber      []string `json:"student_number"`
	AdmissionFeePaid   bool     `json:"admission_fee_paid"`
	InterviewNumber    string   `json:"interview_number"`
	StudyPeriodNumber  string   `json:"study_period_number"`
	CampusClassNumber  string   `json:"campus_class_number"`
	StreamNumber       string   `json:"stream_number"`
	CampusNumber       string   `json:"campus_number"`
	AcademicYearNumber string   `json:"academic_year_number"`
	SchoolNumber       string   `json:"school_number"`
	Items              []string `json:"items"`
}

func (a *AdmissionRequest) Validate() error {
	if len(a.StudentNumber) == 0 {
		return errors.New("at least one student is required")
	}
	if a.CampusNumber == "" {
		return errors.New("campus number is required")
	}
	if a.SchoolNumber == "" {
		return errors.New("school number is required")
	}
	if a.AcademicYearNumber == "" {
		return errors.New("academic year number is required")
	}
	if a.StudyPeriodNumber == "" {
		return errors.New("study period number is required")

	}
	if a.CampusClassNumber == "" {
		return errors.New("campus class number is required")
	}
	if a.StreamNumber == "" {
		return errors.New("stream number is required")
	}
	if len(a.Items) == 0 {
		return errors.New("at least one item is required to be selected for admission")
	}

	return nil
}

func (a *AdmissionRequest) Prepare() {
	a.CampusClassNumber = html.EscapeString(strings.TrimSpace(a.CampusClassNumber))
	a.StreamNumber = html.EscapeString(strings.TrimSpace(a.StreamNumber))
	a.CampusNumber = html.EscapeString(strings.TrimSpace(a.CampusNumber))
	a.AcademicYearNumber = html.EscapeString(strings.TrimSpace(a.AcademicYearNumber))
	a.StudyPeriodNumber = html.EscapeString(strings.TrimSpace(a.StudyPeriodNumber))
	a.SchoolNumber = html.EscapeString(strings.TrimSpace(a.SchoolNumber))

}
