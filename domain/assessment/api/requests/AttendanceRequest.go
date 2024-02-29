package requests

import (
	"errors"
	"time"
)

type AttendanceRequest struct {
	DayOfTheWeek      time.Time `json:"day_of_the_week"`
	CheckIn           time.Time `json:"check_in"`
	CheckOut          time.Time `json:"check_out"`
	Present           bool      `json:"present"`
	StudyPeriodNumber string    `json:"study_period_number"`
	StudentNumber     string    `json:"student_number"`
	ClassNumber       string    `json:"class_number"`
	SchoolNumber      string    `json:"school_number"`
	CampusNumber      string    `json:"campus_number"`
}

func (a *AttendanceRequest) Validate() error {
	if a.DayOfTheWeek == (time.Time{}) {
		return errors.New("day of the week is required")
	}
	if a.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if a.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}

	return nil
}
