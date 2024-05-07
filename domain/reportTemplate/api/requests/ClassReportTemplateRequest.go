package requests

import "errors"

type ClassReportTemplateRequest struct {
	CampusClasses        []string `json:"campus_classes"`
	ReportTemplateNumber string   `json:"report_template_number"`
}

func (classReportTemplateRequest *ClassReportTemplateRequest) Validate() error {

	if classReportTemplateRequest.ReportTemplateNumber == "" {
		return errors.New("choose the template")
	}

	if classReportTemplateRequest.CampusClasses == nil || len(classReportTemplateRequest.CampusClasses) == 0 {
		return errors.New("provide a list of campus classes")
	}

	return nil
}
