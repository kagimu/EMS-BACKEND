package responses

type ItemTransactionHistoryResponse struct {
	Date             string  `json:"date"`
	ItemName         string  `json:"item_name"`
	Amount           float64 `json:"amount"`
	PaymentMethod    string  `json:"payment_method"`
	PaymentReference string  `json:"payment_reference"`
	RecordedBy       string  `json:"recorded_by"`
	StudentDetails   string  `json:"student_details"`
}
