package responses

// PrepaidBalanceResponse is returned by GET prepaid balance.
type PrepaidBalanceResponse struct {
	UserNumber string  `json:"user_number"`
	Balance    float64 `json:"balance"`
}

// PrepaidDepositResponse represents a single deposit record (for list in GET).
type PrepaidDepositResponse struct {
	Number          string  `json:"number"`
	UserNumber      string  `json:"user_number"`
	Amount          float64 `json:"amount"`
	DatePaid        string  `json:"date_paid"`
	Description     string  `json:"description"`
	ClientReference string  `json:"client_reference"`
	CreatedAt       string  `json:"created_at"`
}

// PrepaidBalanceWithDepositsResponse is returned by GET when including deposit history.
type PrepaidBalanceWithDepositsResponse struct {
	UserNumber string                   `json:"user_number"`
	Balance    float64                  `json:"balance"`
	Deposits   []PrepaidDepositResponse `json:"deposits,omitempty"`
}
