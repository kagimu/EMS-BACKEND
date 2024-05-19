package responses

type GroupCalendarRawResponse struct {
	GroupName   string                `json:"group_name"`
	MealCourses []GroupDayMealCourses `json:"meal_courses"`
}
