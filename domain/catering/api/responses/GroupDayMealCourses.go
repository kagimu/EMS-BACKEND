package responses

type GroupDayMealCourses struct {
	Day         string   `json:"day"`
	MealCourses []string `json:"meal_courses"`
	Error       string   `json:"error"`
	Active      bool     `json:"active"`
}

type FeedingGroupCalendarEvent struct {
	CalendarEvents []FeedingGroupCalendarEventItem `json:"calendar_events"`
	CalendarGroups []CalendarGroupItem             `json:"calendar_groups"`
}

type FeedingGroupCalendarEventItem struct {
	CalendarName string `json:"calendar_name"`
	Number       string `json:"number"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
}

type CalendarGroupItem struct {
	GroupName       string               `json:"group_name"`
	Number          string               `json:"number"`
	ActiveCalendars []ActiveCalendarItem `json:"active_calendars"`
}

type ActiveCalendarItem struct {
	FeedingGroupCalendarEventItem
	Active  bool   `json:"active"`
	Message string `json:"message"`
}
