package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"

type GetUserList struct {
	//ID         uint                 `json:"id"`
	Number     string `json:"number"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	LastName   string `json:"last_name"`
	Gender     string `json:"gender"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	/*CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
	DeletedAt  *time.Time           `sql:"index"`*/
	Contacts []models.ContactList `json:"contacts"`
}
