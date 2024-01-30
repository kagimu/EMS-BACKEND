package requests

import "errors"

type ServiceRequest struct {
	Description string `json:"description"`
	Cost        int    `json:"cost"`
	Name        string `json:"name"`
}

func (s *ServiceRequest) Validate() error {
	if s.Description == " " {
		return errors.New("description is required")
	}
	if s.Cost == 0 {
		return errors.New("cost is required")
	}

	return nil
}
