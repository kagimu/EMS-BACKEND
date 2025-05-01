package requests

import (
	"errors"
	"fmt"
)

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
