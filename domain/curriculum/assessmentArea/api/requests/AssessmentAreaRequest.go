package requests

import "errors"

type AssessmentAreaRequest struct {
	Name                  string   `json:"name"`
	CurriculumClassNumber string   `json:"curriculum_class_number"`
	SubjectNumber         []string `json:"subject_numbers"`
}

func (assessmentRequest *AssessmentAreaRequest) Validate() error {

	if len(assessmentRequest.SubjectNumber) == 0 && assessmentRequest.Name == "" {
		return errors.New("both subject number and name can not be empty")
	}

	if len(assessmentRequest.SubjectNumber) > 0 && assessmentRequest.Name != "" {
		return errors.New("either subject number or name has to be empty")
	}

	if assessmentRequest.CurriculumClassNumber == "" {
		return errors.New("curriculum class number can not be empty")
	}

	for _, subjectNumber := range assessmentRequest.SubjectNumber {
		if subjectNumber == "" {
			return errors.New("subject number can not be empty")
		}
	}

	return nil
}
