package requests

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	uuid "github.com/satori/go.uuid"
)

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

func (a *AdmissionLetterTemplateRequest) ToModel() models.AdmissionLetterTemplate {
	return models.AdmissionLetterTemplate{
		Number:       uuid.NewV4().String(),
		SchoolNumber: a.SchoolNumber,
		Phrase:       a.Phrase,
	}
}
