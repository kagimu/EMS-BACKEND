package responses

type TransactionResponse struct {
	ReceiptNumber     string  `json:"receipt_number"`
	PaidBy            string  `json:"paid_by"`
	ReceivedBy        string  `json:"received_by"`
	DateReceived      string  `json:"date_received"`
	DatePaid          string  `json:"date_paid"`
	Amount            float64 `json:"amount"`
	AmountInCredit    float64 `json:"amount_in_credit"`
	Description       string  `json:"description"`
	TransactionNumber string  `json:"transaction_number"`
	Status            string  `json:"status"`
}
