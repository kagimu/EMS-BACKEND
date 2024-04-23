package requests

import "errors"

type InvoiceRequest struct {
	StudentNumber string `json:"student_number"`
	InvoiceNumber string `json:"invoice_number"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	DueDate       string `json:"due_date"`
}

func (i *InvoiceRequest) Validate() error {
	if i.StudentNumber == "" {
		return errors.New("student number is required")
	}

	if i.InvoiceNumber == " " {
		return errors.New("invoice number is required")
	}

	return nil
}
