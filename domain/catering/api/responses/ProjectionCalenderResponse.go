package responses

type ProjectionCalenderResponse struct {
	CalenderItems []ProjectionCalender `json:"calender_items"`
}

type ProjectionCalender struct {
	Day   string                   `json:"day"`
	Date  string                   `json:"date"`
	Meals []ProjectionCalenderMeal `json:"meals"`
}

type ProjectionCalenderMeal struct {
	MealName       string                   `json:"meal_name"`
	Items          []ProjectionCalenderItem `json:"items"`
	NumberOfPeople int                      `json:"number_of_people"`
}

type ProjectionCalenderItem struct {
	Name          string  `json:"name"`
	Quantity      float64 `json:"quantity"`
	UnitOfMeasure string  `json:"unit_of_measure"`
}
