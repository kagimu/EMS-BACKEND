package responses

type StockingResponse struct {
	Number         string  `json:"name"`
	Quantity       int     `json:"quantity"`
	UnitCost       float64 `json:"unit_cost"`
	MenuItemNumber string  `json:"menu_item_number"`
	MenuItemName   string  `json:"menu_item_name"`
	CampusNumber   string  `json:"campus_number"`
	CampusName     string  `json:"campus_name"`
	MenuItemType   string  `json:"menu_item_type"`
}
