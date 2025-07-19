package requests

import (
	"errors"
	"fmt"
	"time"
)

type ClassTeacherCommentRequest struct {
	Comment           string    `json:"comment"`
	CampusClassNumber string    `json:"campus_class_number"`
	StudyPeriodNumber string    `json:"study_period_number"`
	Date              time.Time `json:"date"`
}

func (r *ClassTeacherCommentRequest) Validate() error {
	if r.Comment == "" {
		return errors.New("comment is required")
	}
	if r.CampusClassNumber == "" {
		return errors.New("campus class number is required")
	}
	if r.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}

	return nil
}

type SubjectTeacherCommentRequest struct {
	AssessmentAreaNumber string `json:"assessment_area_number"`
	Comment              string `json:"comment"`
	StudyPeriodNumber    string `json:"study_period_number"`
	Date                 string `json:"date"`
}

func (r *SubjectTeacherCommentRequest) Validate() error {
	if r.AssessmentAreaNumber == "" {
		return errors.New("assessment area number is required")
	}
	if r.Comment == "" {
		return errors.New("comment is required")
	}
	if r.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}
	if r.Date == "" {
		return errors.New("date is required")
	}
	return nil
}

type TeacherCommentRequest struct {
	Comment             string `json:"comment"`
	ReportSummaryNumber string `json:"report_summary_number"`
}

func (r *TeacherCommentRequest) Validate() error {

	if r.Comment == "" {
		return errors.New("comment is required")
	}

	if r.ReportSummaryNumber == "" {
		return errors.New("report summary number is required")
	}

	return nil
}

type LearningAreaReportComment struct {
	Comment            string `json:"comment"`
	LearningAreaNumber string `json:"learning_area_number"`
}

type LearningAreaReportCommentRequest struct {
	ReportSummaryNumber  string                      `json:"report_summary_number"`
	LearningAreaComments []LearningAreaReportComment `json:"learning_area_comments"`
}

func (r *LearningAreaReportCommentRequest) Validate() error {

	if r.ReportSummaryNumber == "" {
		return errors.New("Provide report summary number ")
	}

	if len(r.LearningAreaComments) == 0 {
		return errors.New("Set learning area comments list ")
	}

	for index, comment := range r.LearningAreaComments {

		if comment.LearningAreaNumber == "" {
			return errors.New(fmt.Sprintf("Choose learning area at index %d", index))
		}

		if comment.Comment == "" {
			return errors.New(fmt.Sprintf("Provide comment for learning area at index %d", index))
		}
	}

	return nil
}
