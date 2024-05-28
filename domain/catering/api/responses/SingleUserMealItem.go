package responses

type SingleUserMealItem struct {
	MenuItemNumber string  `json:"menu_item_number"`
	MenuItemName   string  `json:"menu_item_name"`
	Quantity       float64 `json:"quantity"`
	ServedQuantity float64 `json:"served_quantity"`
	UnitOfMeasure  string  `json:"unit_of_measure"`
}
