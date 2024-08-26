package requests

import "errors"

type ClassReportTemplateRequest struct {
	Classes              []string `json:"classes"`
	CampusNumber         string   `json:"campus_number"`
	ReportTemplateNumber string   `json:"report_template_number"`
}

func (classReportTemplateRequest *ClassReportTemplateRequest) Validate() error {

	if classReportTemplateRequest.ReportTemplateNumber == "" {
		return errors.New("choose the template")
	}

	if classReportTemplateRequest.Classes == nil || len(classReportTemplateRequest.Classes) == 0 {
		return errors.New("provide a list of campus classes")
	}

	return nil
}
