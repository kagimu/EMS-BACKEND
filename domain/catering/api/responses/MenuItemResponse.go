package responses

type MenuItemResponse struct {
	Number              string `json:"number"`
	Name                string `json:"name"`
	Code                uint64 `json:"code"`
	Type                string `json:"type"`
	CampusNumber        string `json:"campus_number"`
	CampusName          string `json:"campus_name"`
	Phone               string `json:"phone"`
	Email               string `json:"email"`
	Location            string `json:"location"`
	UnitOfMeasureNumber string `json:"unit_of_measure_number"`
	UnitOfMeasureName   string `json:"unit_of_measure_name"`
	UnitOfMeasureCode   uint64 `json:"unit_of_measure_code"`
	Parts               []Part `json:"parts"`
}

type Part struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Code     uint64  `json:"code"`
	Number   string  `json:"number"`
}
