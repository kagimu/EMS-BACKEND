package requests

import "errors"

type ClassReportTemplateRequest struct {
	CampusClassNumber    string
	ReportTemplateNumber string
}

func (classReportTemplateRequest *ClassReportTemplateRequest) Validate() error {

	if classReportTemplateRequest.ReportTemplateNumber == "" {
		return errors.New("choose the template")
	}

	if classReportTemplateRequest.CampusClassNumber == "" {
		return errors.New("choose the class number")
	}

	return nil
}
