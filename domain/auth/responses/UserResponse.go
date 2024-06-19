package responses

type UserResponse struct {
	Number     string `json:"number"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	LastName   string `json:"last_name"`
	Gender     string `json:"gender"`
	Type       string `json:"type"`
	Status     string `json:"status"`
}
