package requests

import "errors"

type NarrativeRequest struct {
	Description             string `json:"description"`
	StudentNumber           string `json:"student_number"`
	NarrativeScheduleNumber string `json:"narrative_schedule_number"`
	AssessmentAreaNumber    string `json:"assessment_area_number"`
	StudyPeriodNumber       string `json:"study_period_number"`
	CampusClassNumber       string `json:"campus_class_number"`
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

type ScheduledNarrativeRequest struct {
	ScheduledOn          string `json:"scheduled_on"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
	StudyPeriodNumber    string `json:"study_period_number"`
	CampusClassNumber    string `json:"campus_class_number"`
}

func (receiver ScheduledNarrativeRequest) Validate() error {
	if receiver.ScheduledOn == "" {
		return errors.New("Provide the scheduled on date ")
	}

	if receiver.AssessmentAreaNumber == "" {
		return errors.New("Provide the assessment area ")
	}

	if receiver.StudyPeriodNumber == "" {
		return errors.New("Provide the study period ")
	}

	if receiver.CampusClassNumber == "" {
		return errors.New("Provide the campus class ")
	}

	return nil
}
