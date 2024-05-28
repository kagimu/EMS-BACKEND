package responses

type RoleView struct {
	Number      string `json:"number"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RolesList struct {
	Number      string   `json:"number"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}
