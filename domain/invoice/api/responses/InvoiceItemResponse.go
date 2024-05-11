package responses

type InvoiceItemResponse struct {
	ItemNumber string  `json:"item_number"`
	ItemName   string  `json:"item_name"`
	UnitCost   float64 `json:"unit_cost"`
	Quantity   float64 `json:"quantity"`
	TotalCost  float64 `json:"total_cost"`
}
