package requests

import (
	"errors"
)

type InvoiceRequest struct {
	UserNumber  string               `json:"user_number"`
	Description string               `json:"description"`
	DueDate     string               `json:"due_date"`
	Items       []InvoiceItemRequest `json:"items"`
}

type UpdateInvoiceRequest struct {
	InvoiceRequest
	InvoiceNumber string `json:"invoice_number"`
}

func (i *InvoiceRequest) Validate() error {

	if i.UserNumber == "" {
		return errors.New("user number is required")
	}

	if len(i.Items) == 0 {
		return errors.New("add at least one item")
	}

	for _, item := range i.Items {

		if item.Quantity <= 0 {
			return errors.New("quantity must be greater than zero")
		}

		if item.UnitCost < 0 {
			return errors.New("provide unit cost")
		}

	}

	return nil
}

func (i *UpdateInvoiceRequest) Validate() error {

	if i.InvoiceNumber == "" {
		return errors.New("provide invoice number")
	}

	if err := i.InvoiceRequest.Validate(); err != nil {
		return err
	}

	return nil
}
