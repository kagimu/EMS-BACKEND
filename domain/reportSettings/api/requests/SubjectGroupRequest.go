package requests

import (
	"github.com/go-playground/validator/v10"
)

type AssessmentAreaRequest struct {
	AssessmentAreaNumber string `json:"assessment_area_number" validate:"required,uuid4"`
	SubjectNameOnReport  string `json:"subject_name_on_report" validate:"required,string"`
}

type CreateSubjectGroupRequest struct {
	SubjectGroupName       string                  `json:"subject_group_name" validate:"required"`
	ReportSettingsNumber   string                  `json:"report_settings_number" validate:"required,uuid4"`
	ClassNumber            string                  `json:"class_number" validate:"required,uuid4"`
	CampusNumber           string                  `json:"campus_number" validate:"required,uuid4"`
	CurriculumNumber       string                  `json:"curriculum_number" validate:"required,uuid4"`
	CampusCurriculumNumber string                  `json:"campus_curriculum_number" validate:"required,uuid4"`
	AssessmentAreas        []AssessmentAreaRequest `json:"assessmentAreas" validate:"required,dive"`
}

func (r *CreateSubjectGroupRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type AddAssessmentAreaSubjectGroupRequest struct {
	SubjectGroupNumber string                  `json:"subject_group_number" validate:"required,uuid4"`
	AssessmentAreas    []AssessmentAreaRequest `json:"assessmentAreas" validate:"required,dive"`
}

func (r *AddAssessmentAreaSubjectGroupRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
