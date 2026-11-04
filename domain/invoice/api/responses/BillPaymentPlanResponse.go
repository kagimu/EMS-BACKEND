package responses

type BillPaymentPlanInstallmentPaymentResponse struct {
	Number     string  `json:"number"`
	PaidAmount float64 `json:"paid_amount"`
	PaidOn     string  `json:"paid_on"`
}

type BillPaymentPlanInstallmentResponse struct {
	Number     string                                      `json:"number"`
	DueDate    string                                      `json:"due_date"`
	Amount     float64                                     `json:"amount"`
	PaidAmount float64                                     `json:"paid_amount"`
	Status     string                                      `json:"status"`
	Sequence   int                                         `json:"sequence"`
	Payments   []BillPaymentPlanInstallmentPaymentResponse `json:"bill_payment_plan_installment_payments"`
}

type BillPaymentPlanResponse struct {
	Number               string                               `json:"number"`
	PromissoryNoteNumber string                               `json:"promissory_note_number"`
	StudentNumber        string                               `json:"student_number"`
	StudyPeriodNumber    string                               `json:"study_period_number"`
	TotalAmount          float64                              `json:"total_amount"`
	PaidAmount           float64                              `json:"paid_amount"`
	StartDate            string                               `json:"start_date"`
	Notes                string                               `json:"notes"`
	Status               string                               `json:"status"`
	Installments         []BillPaymentPlanInstallmentResponse `json:"installments"`
}
