package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"

type GetUserList struct {
	//ID         uint                 `json:"id"`
	Number     string               `json:"number"`
	FirstName  string               `json:"first_name"`
	SecondName string               `json:"second_name"`
	LastName   string               `json:"last_name"`
	Gender     string               `json:"gender"`
	Type       string               `json:"type"`
	Status     string               `json:"status"`
	Contacts   []models.ContactList `json:"contacts"`
}

type UserView struct {
	Number      string               `json:"number"`
	FirstName   string               `json:"first_name"`
	SecondName  string               `json:"second_name"`
	LastName    string               `json:"last_name"`
	Gender      string               `json:"gender"`
	Type        string               `json:"type"`
	Status      string               `json:"status"`
	Contacts    []models.ContactList `json:"contacts"`
	Roles       []RoleView           `json:"roles"`
	Permissions []string             `json:"permissions"`
}

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
