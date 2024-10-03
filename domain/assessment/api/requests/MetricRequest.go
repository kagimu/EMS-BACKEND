package requests

import "errors"

type MetricRequest struct {
	Name             string `json:"name"`
	CurriculumNumber string `json:"curriculum_number"`
}

func (m *MetricRequest) Validate() error {
	if m.Name == "" {
		return errors.New("name is required")
	}
	if m.CurriculumNumber == "" {
		return errors.New("curriculum number is required")
	}
	return nil
}
