package requests

import "errors"

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
