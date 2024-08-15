package requests

import "errors"

type StudyPeriodDelRequest struct {
	Number string `json:"number"`
}

type StudyPeriodInfo struct {
	Code                 string `json:"code"`
	StudyPeriodName      string `json:"study_period_name"`
	AcademicYearNumber   string `json:"academic_year_number"`
	StudyPeriodStartDate string `json:"study_period_start_date"`
	StudyPeriodEndDate   string `json:"study_period_end_date"`
	StudyPeriodNumber    string `json:"study_period_number"`
}

type StudyPeriodNameRequest struct {
	SchoolNumber string `json:"school_number"`
	Name         string `json:"name"`
}

type UpdateStudyPeriodNameRequest struct {
	StudyPeriodNameNumber string `json:"study_period_name_number"`
	StudyPeriodNameRequest
}

func (studyPeriodNameRequest *StudyPeriodNameRequest) Validate() error {

	if studyPeriodNameRequest.SchoolNumber == "" {
		return errors.New("Choose school ")
	}

	if studyPeriodNameRequest.Name == "" {
		return errors.New("Provide the name ")
	}

	if len(studyPeriodNameRequest.Name) > 60 {
		return errors.New("Study period name too long, make at most 60 characters ")
	}

	return nil
}

func (updateStudyPeriodNameRequest *UpdateStudyPeriodNameRequest) Validate() error {
	if updateStudyPeriodNameRequest.StudyPeriodNameNumber == "" {
		return errors.New("Choose the name to edit ")
	}

	if err := updateStudyPeriodNameRequest.StudyPeriodNameRequest.Validate(); err != nil {
		return err
	}

	return nil
}
