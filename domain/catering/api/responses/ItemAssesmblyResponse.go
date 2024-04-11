package responses

type ItemAssemblyResponse struct {
	Number              string  `json:"number"`
	PrimaryItemNumber   string  `json:"primary_item_number"`
	SecondaryItemNumber string  `json:"secondary_item_number"`
	Quantity            float64 `json:"quantity"`
	PrimaryItemName     string  `json:"primary_item_name"`
	SecondaryItemName   string  `json:"secondary_item_name"`
}
