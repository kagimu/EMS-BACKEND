package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/responses"

type StaffResponse struct {
	Number string `json:"number"`
	responses.UserResponse
	Classes []string `json:"classes"`
	Streams []string `json:"streams"`
}
