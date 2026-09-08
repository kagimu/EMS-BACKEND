package requests

import "errors"

type ItemPayment struct {
	InvoiceItemNumber string  `json:"invoice_item_number"`
	Amount            float64 `json:"amount"`
}

type InvoicePaymentRequest struct {
	InvoiceNumber        string        `json:"invoice_number"`
	Amount               float64       `json:"amount"`
	Date                 string        `json:"date"`
	Credits              []string      `json:"credits"`
	ItemPayments         []ItemPayment `json:"item_payments,omitempty"`
	PaymentMethod        string        `json:"payment_method,omitempty"`
	PaymentReference     string        `json:"payment_reference,omitempty"`
	DepositAccountNumber string        `json:"deposit_account_number,omitempty"`
}

type AccountTopUpRequest struct {
	UserNumber           string  `json:"user_number"`
	Amount               float64 `json:"amount"`
	Description          *string `json:"description,omitempty"`
	DatePaid             string  `json:"date_paid,omitempty"`
	ClientReference      *string `json:"client_reference,omitempty"`
	DepositAccountNumber string  `json:"deposit_account_number,omitempty"`
}

func (r *AccountTopUpRequest) Validate() error {
	if r.UserNumber == "" {
		return errors.New("choose the user paying")
	}
	if r.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if r.DatePaid == "" {
		return errors.New("choose the date paid")
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
