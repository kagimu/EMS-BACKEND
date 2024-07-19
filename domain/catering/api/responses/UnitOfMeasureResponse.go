package responses

type UnitOfMeasureResponse struct {
	UnitOfMeasureNumber string `json:"unit_of_measure_number"`
	Name                string `json:"name"`
	Abbreviation        string `json:"abbreviation"`
	Code                uint64 `json:"code"`
	Status              string `json:"status"`
	SchoolNumber        string `json:"school_number"`
}
