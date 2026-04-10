package responses

type BillPaymentPlanInstallmentResponse struct {
	Number   string  `json:"number"`
	DueDate  string  `json:"due_date"`
	Amount   float64 `json:"amount"`
	PaidAmount float64 `json:"paid_amount"`
	Status   string  `json:"status"`
	Sequence int     `json:"sequence"`
}

type BillPaymentPlanResponse struct {
	Number            string                             `json:"number"`
	BillableNumber    string                             `json:"billable_number"`
	StudentNumber     string                             `json:"student_number"`
	StudyPeriodNumber string                             `json:"study_period_number"`
	TotalAmount       float64                            `json:"total_amount"`
	StartDate         string                             `json:"start_date"`
	Notes             string                             `json:"notes"`
	Status            string                             `json:"status"`
	Installments      []BillPaymentPlanInstallmentResponse `json:"installments"`
}

