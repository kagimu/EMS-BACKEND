package requests

type CourseRequest struct {
	Name string `json:"name"`
	Code int    `json:"code"`
}
