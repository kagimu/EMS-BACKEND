package responses

type ProjectionCalenderResponse struct {
	CalenderItems []ProjectionCalender `json:"calender_items"`
	ItemsSummary  []ProjectionMenuItem `json:"items_summary"`
	TotalCost     float64              `json:"total_cost"`
}

type ProjectionCalender struct {
	Day   string                   `json:"day"`
	Date  string                   `json:"date"`
	Meals []ProjectionCalenderMeal `json:"meals"`
}

type ProjectionCalenderMeal struct {
	MealName       string               `json:"meal_name"`
	Items          []ProjectionMenuItem `json:"items"`
	NumberOfPeople int                  `json:"number_of_people"`
}

type ProjectionMenuItem struct {
	Name           string  `json:"name"`
	Quantity       float64 `json:"quantity"`
	UnitOfMeasure  string  `json:"unit_of_measure"`
	TotalCost      float64 `json:"total_cost"`
	MenuItemNumber string  `json:"menu_item_number"`
}
