package requests

type UserDataRequest struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	LastName   string `json:"last_name"`
	Gender     string `json:"gender"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}
