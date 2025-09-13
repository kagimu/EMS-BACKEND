package requests

import "errors"

type NarrativeRequest struct {
	Description             string `json:"description"`
	StudentNumber           string `json:"student_number"`
	NarrativeScheduleNumber string `json:"narrative_schedule_number"`
	AssessmentAreaNumber    string `json:"assessment_area_number"`
	StudyPeriodNumber       string `json:"study_period_number"`
	CampusClassNumber       string `json:"campus_class_number"`
	ReferenceDate           string `json:"reference_date"`
	ReferenceTime           string `json:"reference_time"`
	CurrentStatus           string `json:"current_status"`
	Status                  string `json:"status"`
	SuggestedSolution       string `json:"suggested_solution"`
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

	if receiver.Status == "" {
		return errors.New("Choose status ")
	}

	if receiver.CurrentStatus == "" {
		return errors.New("Choose current status ")
	}

	if receiver.ReferenceDate == "" {
		return errors.New("Choose the date ")
	}

	if receiver.ReferenceTime == "" {
		return errors.New("Choose the time ")
	}

	//if receiver.NarrativeScheduleNumber == "" {
	//	return errors.New("Please Choose the schedule to enter narrative scores ")
	//}

	return nil
}

type ScheduledNarrativeRequest struct {
	ScheduledOn          []string `json:"scheduled_on"`
	AssessmentAreaNumber string   `json:"assessment_area_number"`
	StudyPeriodNumber    string   `json:"study_period_number"`
	CampusNumber         string   `json:"campus_number"`
	ClassNumber          string   `json:"class_number"`
	CurriculumNumber     string   `json:"curriculum_number"`
	Description          string   `json:"description"`
}

func (receiver ScheduledNarrativeRequest) Validate() error {
	if len(receiver.ScheduledOn) == 0 {
		return errors.New("Provide the scheduled on dates ")
	}

	if receiver.AssessmentAreaNumber == "" {
		return errors.New("Provide the assessment area ")
	}

	if receiver.StudyPeriodNumber == "" {
		return errors.New("Provide the study period ")
	}

	if receiver.CampusNumber == "" {
		return errors.New("Choose campus ")
	}

	if receiver.CurriculumNumber == "" {
		return errors.New("Choose curriculum ")
	}

	if receiver.ClassNumber == "" {
		return errors.New("Choose class ")
	}

	if receiver.Description == "" {
		return errors.New("Provide the description ")
	}

	if len(receiver.Description) > 2500 {
		return errors.New("Provide the description length less than 2500 ")
	}

	return nil
}
