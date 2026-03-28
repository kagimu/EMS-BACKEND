package requests

import (
	"errors"
)

type SubjectTeacherReportRequest struct {
	SubjectTeacherName   string `json:"subject_teacher_name"`
	Status               string `json:"status"`
	ClassTeacherComment  string `json:"class_teacher_comment"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
	StudyPeriodNumber    string `json:"study_period_number"`
	Stage                string `json:"stage"`
}

func (r *SubjectTeacherReportRequest) Validate() error {
	if r.SubjectTeacherName == "" {
		return errors.New("subject teacher name is required")
	}
	if r.AssessmentAreaNumber == "" {
		return errors.New("assessment area number is required")
	}
	if r.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}

	return nil
}

type UpdateSubjectTeacherReportRequest struct {
	SubjectTeacherName   string `json:"subject_teacher_name"`
	Status               string `json:"status"`
	ClassTeacherComment  string `json:"class_teacher_comment"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
	StudyPeriodNumber    string `json:"study_period_number"`
	Stage                string `json:"stage"`
}

func (r *UpdateSubjectTeacherReportRequest) Validate() error {
	// Add update-specific validation if needed
	return nil
}

type ForwardSubjectTeacherReportRequest struct {
	SubjectTeacherReportNumber string `json:"subject_teacher_report_number" binding:"required"`
	NextStage                  string `json:"next_stage" binding:"required"`
	Comment                    string `json:"comment" binding:"required"`
}

func (r *ForwardSubjectTeacherReportRequest) Validate() error {
	if r.SubjectTeacherReportNumber == "" {
		return errors.New("subject teacher report number is required")
	}
	if r.NextStage == "" {
		return errors.New("next stage is required")
	}
	if r.Comment == "" {
		return errors.New("comment is required")
	}
	return nil
}
