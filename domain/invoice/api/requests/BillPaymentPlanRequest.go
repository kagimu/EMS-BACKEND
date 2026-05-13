package requests

import "errors"

type BillPaymentPlanInstallmentRequest struct {
	DueDate string  `json:"due_date"`
	Amount  float64 `json:"amount"`
}

type BillPaymentPlanRequest struct {
	PromissoryNoteNumber string                              `json:"promissory_note_number"`
	StudentNumber        string                              `json:"student_number"`
	StartDate            string                              `json:"start_date"`
	Notes                string                              `json:"notes"`
	Installments         []BillPaymentPlanInstallmentRequest `json:"installments"`
}

type UpdateBillPaymentPlanRequest struct {
	BillPaymentPlanRequest
	Status string `json:"status"`
}

func (r *BillPaymentPlanRequest) Validate() error {
	if r.PromissoryNoteNumber == "" {
		return errors.New("choose promissory note")
	}
	if r.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if r.StartDate == "" {
		return errors.New("start date is required")
	}
	if len(r.Installments) == 0 {
		return errors.New("at least one installment is required")
	}

	for _, installment := range r.Installments {
		if installment.DueDate == "" {
			return errors.New("installment due date is required")
		}
		if installment.Amount <= 0 {
			return errors.New("installment amount must be greater than zero")
		}
	}

	return nil
}

func (r *UpdateBillPaymentPlanRequest) Validate() error {
	return r.BillPaymentPlanRequest.Validate()
}

type PayInstallmentRequest struct {
	InstallmentNumber string  `json:"installment_number"`
	PaidAmount        float64 `json:"paid_amount"`
	PaidOn            string  `json:"paid_on"`
}

func (r *PayInstallmentRequest) Validate() error {
	if r.InstallmentNumber == "" {
		return errors.New("installment number is required")
	}
	if r.PaidAmount <= 0 {
		return errors.New("paid amount must be greater than zero")
	}
	if r.PaidOn == "" {
		return errors.New("payment date is required")
	}
	return nil
}
