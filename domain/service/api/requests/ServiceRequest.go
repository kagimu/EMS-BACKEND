package requests

import "errors"

type ServiceRequest struct {
	Description   string   `json:"description"`
	Cost          int      `json:"cost"`
	Name          string   `json:"name"`
	Mandatory     bool     `json:"mandatory"`
	CampusNumbers []string `json:"campus_number"`
	PaymentMode   string   `json:"payment_mode"`
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
