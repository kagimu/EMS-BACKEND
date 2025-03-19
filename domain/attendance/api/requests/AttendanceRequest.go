package requests

import (
	"errors"
)

type CheckoutRequest struct {
	Number       string `json:"number"`
	CheckoutTime string `json:"checkout_time"`
	Comment      string `json:"comment"`
}

func (c *CheckoutRequest) ValidateCheckout() error {
	if c.Number == "" {
		return errors.New("please input checkout number")
	}

	if c.CheckoutTime == "" {
		return errors.New("please input checkout time")
	}

	return nil
}

type AttendanceRequest struct {
	Date              string `json:"date"`
	CheckInTime       string `json:"check_in_time"`
	CheckOutTime      string `json:"check_out_time"`
	Present           string `json:"present"`
	UserNumber        string `json:"user_number"`
	CampusClassNumber string `json:"campus_class_number"`
	SchoolNumber      string `json:"school_number"`
	CampusNumber      string `json:"campus_number"`
	AbsenceReason     string `json:"absence_reason"`
	StudentNumber     string `json:"student_number"`
	StreamNumber      string `json:"stream_number"`
	StudyPeriodNumber string `json:"study_period_number"`
}

func (a *AttendanceRequest) Validate() error {
	if a.Date == "" {
		return errors.New("date of attendance is required")
	}
	if a.UserNumber == "" {
		return errors.New("please select the student/staff you are capturing attendance for ")
	}

	if a.Present == "" {
		return errors.New("please indicate if the student/staff is present or absent")
	}

	if a.Present == "NO" && a.AbsenceReason == "" {
		return errors.New("please provide a reason for the absence")
	}

	if a.Present == "YES" && a.CheckInTime == "" {
		return errors.New("please provide the check in time")
	}

	return nil
}
