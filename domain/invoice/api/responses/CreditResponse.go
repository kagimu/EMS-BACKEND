package responses

type CreditResponse struct {
	Number          string  `json:"number"`
	Date            string  `json:"date"`
	Amount          float64 `json:"amount"`
	Balance         float64 `json:"balance"`
	Status          string  `json:"status"`
	ClientReference string  `json:"client_reference"`
	Description     string  `json:"description"`
	SchoolNumber    string  `json:"school_number"`
	CampusNumber    string  `json:"campus_number"`
	UserNumber      string  `json:"user_number"`
	AddedByNumber   string  `json:"added_by_number"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}
