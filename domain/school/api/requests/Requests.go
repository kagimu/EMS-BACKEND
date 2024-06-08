package requests

import "errors"

type SchoolInfo struct {
	Name             string   `json:"name"`
	EntityNumber     string   `json:"entity_number"`
	CommencementDate string   `json:"commencement_date"`
	SchoolNumber     string   `json:"school_number"`
	Curriculum       []string `json:"curriculum"`
}

func (r *SchoolInfo) Validate() error {
	if r.Name == "" {
		return errors.New("school name is required")
	}
	if r.EntityNumber == "" {
		return errors.New("entity number is required")
	}
	if r.CommencementDate == "" {
		return errors.New("commencement date is required")
	}
	if len(r.Curriculum) == 0 {
		return errors.New("at least one curriculum is required")
	}

	return nil
}
