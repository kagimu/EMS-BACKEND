package requests

import "errors"

type NarrativeRequest struct {
	Description          string `json:"description"`
	StudentNumber        string `json:"student_number"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
	StudyPeriodNumber    string `json:"study_period_number"`
	CampusClassNumber    string `json:"campus_class_number"`
}

func (receiver NarrativeRequest) Validate() error {
	if receiver.Description == "" {
		return errors.New("Provide the narrative description ")
	}

	if receiver.StudentNumber == "" {
		return errors.New("Provide the student number ")
	}

	if receiver.AssessmentAreaNumber == "" {
		return errors.New("Assessment area is required ")
	}

	if receiver.StudyPeriodNumber == "" {
		return errors.New("Choose study period ")
	}

	if receiver.CampusClassNumber == "" {
		return errors.New("Choose the campus class ")
	}

	return nil
}
