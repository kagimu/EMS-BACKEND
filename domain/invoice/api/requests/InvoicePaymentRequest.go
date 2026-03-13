package requests

import "errors"

type InvoicePaymentRequest struct {
	InvoiceNumber string   `json:"invoice_number"`
	Amount        float64  `json:"amount"`
	Date          string   `json:"date"`
	Credits       []string `json:"credits"`
}

// PrepaidDepositRequest is the body for depositing money into a user's prepaid balance.
type PrepaidDepositRequest struct {
	UserNumber      string  `json:"user_number"`
	Amount          float64 `json:"amount"`
	Description     *string `json:"description"`
	DatePaid        string  `json:"date_paid"`
	ClientReference *string `json:"client_reference"`
}

func (r *PrepaidDepositRequest) Validate() error {
	if r.UserNumber == "" {
		return errors.New("user_number is required")
	}
	if r.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if r.DatePaid == "" {
		return errors.New("date_paid is required")
	}
	return nil
}

func (i *InvoicePaymentRequest) Validate() error {

	if len(i.Credits) == 0 {
		if i.Amount <= 0 {
			return errors.New("amount should be greater than zero")
		}
		if i.Date == "" {
			return errors.New("date when paid is required")
		}
	} else {
		for _, creditNumber := range i.Credits {
			if creditNumber == "" {
				return errors.New("You have empty credit ID ")
			}
		}
	}

	if i.InvoiceNumber == "" {
		return errors.New("choose the promissory note they are paying against")
	}

	return nil
}
