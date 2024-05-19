package responses

type FeedingGroupMealCourseCalendar struct {
	Headers      []string                   `json:"headers"`
	FeedingGroup []GroupCalendarRawResponse `json:"feeding_group"`
}
