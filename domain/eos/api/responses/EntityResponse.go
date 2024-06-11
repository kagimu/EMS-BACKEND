package responses

type EntityResponse struct {
	Name    string `json:"name"`
	Number  string `json:"number"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}
