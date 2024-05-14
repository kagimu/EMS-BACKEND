package responses

type MealCourseItemResponse struct {
	Name    string `json:"name"`
	AddedBy string `json:"added_by"`
	Code    int    `json:"code"`
}
