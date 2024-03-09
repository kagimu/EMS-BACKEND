package requests

import "errors"

type AssessmentAreaRequest struct {
	Name                  string `json:"name"`
	CurriculumClassNumber string `json:"curriculum_class_number"`
	SubjectNumber         string `json:"subject_number"`
}

func (a *AssessmentAreaRequest) Validate() error {
	if a.CurriculumClassNumber == "" {
		return errors.New("curriculum class number is required")
	}
	if a.SubjectNumber == "" {
		return errors.New("subject number is required")
	}
	if a.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
