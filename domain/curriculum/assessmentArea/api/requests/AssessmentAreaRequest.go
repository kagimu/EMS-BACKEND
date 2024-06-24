package requests

import "errors"

type AssessmentAreaRequest struct {
	CurriculumClassNumber string   `json:"curriculum_class_number"`
	Subjects              []string `json:"subjects"`
}

func (assessmentRequest *AssessmentAreaRequest) Validate() error {

	if len(assessmentRequest.Subjects) == 0 {
		return errors.New("choose subject numbers or custom names")
	}

	if assessmentRequest.CurriculumClassNumber == "" {
		return errors.New("curriculum class number can not be empty")
	}

	for _, subjectNumber := range assessmentRequest.Subjects {
		if subjectNumber == "" {
			return errors.New("subject number can not be empty")
		}
	}

	return nil
}

type UpdateAssessmentAreaRequest struct {
	Name                 string `json:"name"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
}

func (updateAssessmentRequest *UpdateAssessmentAreaRequest) Validate() error {

	if updateAssessmentRequest.AssessmentAreaNumber == "" {
		return errors.New("provide assessment area number being edited")
	}

	if updateAssessmentRequest.Name == "" {
		return errors.New("provide the assessment area updated name")
	}

	return nil
}
