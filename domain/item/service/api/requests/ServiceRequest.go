package requests

import (
	"errors"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ServiceRequest struct {
	Description     string   `json:"description"`
	Cost            int64    `json:"cost"`
	Name            string   `json:"name"`
	Mandatory       bool     `json:"mandatory"`
	SerialCode      string   `json:"serial_code"`
	CampusNumbers   []string `json:"campus_numbers"`
	OneTimePurchase bool     `json:"one_time_purchase"`
	RepeatCycle     string   `json:"repeat_cycle"`
	PaymentMode     string   `json:"payment_mode"`
}

func (s *ServiceRequest) Validate() error {
	if s.Description == " " {
		return errors.New("description is required")
	}
	if s.Cost == 0 {
		return errors.New("cost is required")
	}
	if s.Name == "" {
		return errors.New("name is required")
	}
	if len(s.CampusNumbers) == 0 {
		return errors.New("at least one campus is required")
	}

	return nil
}

func (s *ServiceRequest) Prepare() {
	s.Name = cases.Title(language.Und).String(s.Name)
}
