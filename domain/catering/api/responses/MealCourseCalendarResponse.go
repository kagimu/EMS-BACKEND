package responses

type MealCourseCalendarResponse struct {
	Headers             []string              `json:"headers"`
	CalendarRawResponse []CalendarRowResponse `json:"calendar_raw_response"`
}
