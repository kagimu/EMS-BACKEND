package responses

type GroupDayMealCourses struct {
	Day         string   `json:"day"`
	MealCourses []string `json:"meal_courses"`
	Error       string   `json:"error"`
}
