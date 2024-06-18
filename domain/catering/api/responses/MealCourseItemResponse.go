package responses

type MealCourseItemResponse struct {
	Name    string `json:"name"`
	AddedBy string `json:"added_by"`
	Code    int    `json:"code"`
	Number  string `json:"number"`
	Status  string `json:"status"`
}
