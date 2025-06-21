package requests

import (
	"errors"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ServiceRequest struct {
	Description     string   `json:"description"`
	Cost            float64  `json:"cost"`
	Name            string   `json:"name"`
	SerialCode      string   `json:"serial_code"`
	Mandatory       bool     `json:"mandatory"`
	CampusNumbers   []string `json:"campus_numbers"`
	OneTimePurchase bool     `json:"one_time_purchase"`
	RepeatCycle     string   `json:"repeat_cycle"`
	PaymentMode     string   `json:"payment_mode"`
	SchoolNumber    string   `json:"school_number"`
}

func (s *ServiceRequest) Validate() error {
	if s.Description == " " {
		return errors.New("description is required")
	}
	if s.Cost < 0 {
		return errors.New("cost can not be negative")
	}
	if s.Name == "" {
		return errors.New("name is required")
	}
	if len(s.CampusNumbers) == 0 {
		return errors.New("at least one campus is required")
	}
	if s.SchoolNumber == "" {
		return errors.New("school number is required")
	}

	return nil
}

func (s *ServiceRequest) Prepare() {
	s.Name = cases.Title(language.Und).String(s.Name)
}
