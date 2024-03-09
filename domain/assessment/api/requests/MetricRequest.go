package requests

import "errors"

type MetricRequest struct {
	Name             string `json:"name"`
	CurriculumNumber string `json:"curriculum_number"`
	ExamType         string `json:"exam_type"`
}

func (m *MetricRequest) Validate() error {
	if m.Name == "" {
		return errors.New("name is required")
	}
	if m.CurriculumNumber == "" {
		return errors.New("curriculum number is required")
	}
	if m.ExamType == "" {
		return errors.New("exam_type is required")
	}
	return nil
}
