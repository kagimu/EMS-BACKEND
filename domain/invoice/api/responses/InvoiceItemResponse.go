package responses

type InvoiceItemResponse struct {
	ItemNumber string  `json:"item_number"`
	ItemName   string  `json:"item_name"`
	Quantity   float64 `json:"quantity"`
	UnitCost   float64 `json:"unit_cost"`
	TotalCost  float64 `json:"total_cost"`
	PaidAmount float64 `json:"paid_amount"`
	Balance    float64 `json:"balance"`
}
