package requests

import (
	"errors"
)

type InvoiceRequest struct {
	UserNumber  string               `json:"student_number"`
	Description string               `json:"description"`
	DueDate     string               `json:"due_date"`
	Items       []InvoiceItemRequest `json:"items"`
}

func (i *InvoiceRequest) Validate() error {

	if i.UserNumber == "" {
		return errors.New("student number is required")
	}

	if len(i.Items) == 0 {
		return errors.New("add at least one item")
	}

	return nil
}
