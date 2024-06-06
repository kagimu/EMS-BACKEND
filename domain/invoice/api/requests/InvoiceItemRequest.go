package requests

type InvoiceItemRequest struct {
	ItemNumber string  `json:"item_number"`
	Quantity   float64 `json:"quantity"`
	UnitCost   float64 `json:"unit_cost"`
}
