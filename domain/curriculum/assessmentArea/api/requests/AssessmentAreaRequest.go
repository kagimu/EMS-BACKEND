package requests

import "errors"

type AssessmentAreaRequest struct {
	Name                  string `json:"name"`
	CurriculumClassNumber string `json:"curriculum_class_number"`
	SubjectNumber         string `json:"subject_number"`
}

func (a *AssessmentAreaRequest) Validate() error {
	if a.SubjectNumber == "" && a.Name == "" {
		return errors.New("both subject number and name can not be empty")
	}
	if a.SubjectNumber != "" && a.Name != "" {
		return errors.New("either subject number or name has to be empty")
	}
	if a.CurriculumClassNumber == "" {
		return errors.New("curriculum class number can not be empty")
	}

	return nil
}
