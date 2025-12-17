package responses

type FinancialSummaryResponse struct {
	TotalInvoices      int64      `json:"total_invoices"`
	TotalInvoiceAmount float64    `json:"total_invoice_amount"`
	TotalPayments      int64      `json:"total_payments"`
	TotalPaidAmount    float64    `json:"total_paid_amount"`
	TotalBalanceAmount float64    `json:"total_balance_amount"`
	TotalCreditAmount  float64    `json:"total_credit_amount"`
	MostRecent         MostRecent `json:"most_recent"`
}

type MostRecent struct {
	InvoiceAmount float64 `json:"invoice_amount"`
	InvoiceDate   string  `json:"invoice_date"`
	PaidAmount    float64 `json:"paid_amount"`
	PaidDate      string  `json:"paid_date"`
	DateReceived  string  `json:"date_received"`
}
