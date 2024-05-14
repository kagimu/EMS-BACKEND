package responses

type CalendarRowResponse struct {
	CourseName string             `json:"course_name"`
	DayItems   []DayItemsResponse `json:"day_items"`
}
