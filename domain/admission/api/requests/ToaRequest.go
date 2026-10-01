package requests

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	uuid "github.com/satori/go.uuid"
)

type ToaRequest struct {
	SchoolNumber string `json:"school_number" binding:"required"`
	Content      string `json:"content" binding:"required"`
}

func (r *ToaRequest) Validate() error {
	if r.SchoolNumber == "" {
		return errors.New("school number is required")
	}
	if r.Content == "" {
		return errors.New("content is required")
	}
	return nil
}

func (r *ToaRequest) ToModel() models.TermsOfAdmission {
	return models.TermsOfAdmission{
		Number:       uuid.NewV4().String(),
		SchoolNumber: r.SchoolNumber,
		Content:      r.Content,
	}
}
