package responses

type EntityResponse struct {
	Name      string `json:"name"`
	Number    string `json:"number"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Facebook  string `json:"facebook"`
	Instagram string `json:"instagram"`
	Twitter   string `json:"twitter"`
	Website   string `json:"website"`
	Whatsapp  string `json:"whatsapp"`
	Linkedin  string `json:"linkedin"`
}
