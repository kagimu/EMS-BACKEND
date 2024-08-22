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
		return errors.New("school name is required")
	}
	if r.EntityNumber == "" {
		return errors.New("entity number is required")
	}
	if r.CommencementDate == "" {
		return errors.New("commencement date is required")
	}
	if len(r.Curriculum) == 0 {
		return errors.New("at least one curriculum is required")
	}

	return nil
}

type ModeOfStudyRequest struct {
	SchoolNumber string `json:"school_number"`
	Name         string `json:"name"`
}

type UpdateModeOfStudyRequest struct {
	ModeOfStudyNumber string `json:"mode_of_study_number"`
	Status            string `json:"status"`
	ModeOfStudyRequest
}

func (modeOfStudyRequest *ModeOfStudyRequest) Validate() error {
	if modeOfStudyRequest.SchoolNumber == "" {
		return errors.New("Choose school ")
	}

	if modeOfStudyRequest.Name == "" {
		return errors.New("Enter mode name I.E Boarding, Half Day ")
	}

	return nil
}

func (modeOfStudyRequest *UpdateModeOfStudyRequest) Validate() error {
	if modeOfStudyRequest.ModeOfStudyNumber == "" {
		return errors.New("Choose mode of study being updated ")
	}

	return modeOfStudyRequest.ModeOfStudyRequest.Validate()
}
