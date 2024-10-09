package responses

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/api/responses"
)

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
	Roles       []responses.RoleView `json:"roles"`
	Permissions []string             `json:"permissions"`
	Scopes      []string             `json:"scopes"`
}
