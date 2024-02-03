package requests

import "errors"

type InvoiceRequest struct {
	Status           string   `json:"status"`
	StudentNumber    string   `json:"student_number"`
	PaymentReference string   `json:"payment_reference"`
	ServiceNumber    []string `json:"service_number"`
}

func (i *InvoiceRequest) Validate() error {
	if i.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if len(i.ServiceNumber) == 0 {
		return errors.New("service number is required")
	}

	return nil
}
