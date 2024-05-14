package responses

type MealCourseResponse struct {
	Count        int                      `json:"count"`
	SchoolNumber string                   `json:"school_number"`
	Courses      []MealCourseItemResponse `json:"courses"`
}
