package requests

import "errors"

type AdmissionLetterTemplateRequest struct {
	SchoolNumber string `json:"school_number" binding:"required"`
	Phrase       string `json:"phrase" binding:"required"`
}

func (a *AdmissionLetterTemplateRequest) Validate() error {
	if a.SchoolNumber == "" {
		return errors.New("school number is required")
	}
	if a.Phrase == "" {
		return errors.New("phrase is required")
	}
	return nil
}
