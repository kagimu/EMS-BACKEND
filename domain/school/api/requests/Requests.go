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
		return errors.New("School name is required")
	}
	if r.EntityNumber == "" {
		return errors.New("Entity number is required")
	}
	if r.CommencementDate == "" {
		return errors.New("Commencement date is required")
	}
	if len(r.Curriculum) == 0 {
		return errors.New("At least one curriculum is required")
	}

	return nil
}
