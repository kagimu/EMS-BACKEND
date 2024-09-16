package responses

type StatementResponse struct {
	Amount       float64 `json:"amount"`
	Balance      float64 `json:"balance"`
	Type         string  `json:"type"`
	DateReceived string  `json:"date_received"`
	ReceivedBy   string  `json:"received_by"`
	Description  string  `json:"description"`
}
