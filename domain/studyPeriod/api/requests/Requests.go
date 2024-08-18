package requests

import "errors"

type StudyPeriodDelRequest struct {
	Number string `json:"number"`
}

type StudyPeriodRequest struct {
	Code                  int64    `json:"code"`
	AcademicYearNumber    string   `json:"academic_year_number"`
	StudyPeriodStartDate  string   `json:"study_period_start_date"`
	StudyPeriodEndDate    string   `json:"study_period_end_date"`
	StudyPeriodNameNumber string   `json:"study_period_name_number"`
	Classes               []string `json:"classes"`
}

func (studyPeriodRequest *StudyPeriodRequest) Validate() error {

	if studyPeriodRequest.Code < 1 {
		return errors.New("Provide the code ")
	}

	if studyPeriodRequest.StudyPeriodStartDate == "" {
		return errors.New("Choose start date ")
	}

	if studyPeriodRequest.StudyPeriodEndDate == "" {
		return errors.New("Choose end date ")
	}

	if studyPeriodRequest.StudyPeriodNameNumber == "" {
		return errors.New("Choose study period name ")
	}

	if studyPeriodRequest.StudyPeriodNameNumber == "" {
		return errors.New("Choose the academic year ")
	}

	if len(studyPeriodRequest.Classes) == 0 {
		return errors.New("Choose the classes for this study period ")
	}

	return nil
}

type UpdateStudyPeriodRequest struct {
	StudyPeriodRequest
	StudyPeriodNumber string `json:"study_period_number"`
}

func (updateStudyPeriodRequest *UpdateStudyPeriodRequest) Validate() error {
	if updateStudyPeriodRequest.StudyPeriodNumber == "" {
		return errors.New("Choose study period to update ")
	}

	err := updateStudyPeriodRequest.StudyPeriodRequest.Validate()
	if err != nil {
		return err
	}

	return nil
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
