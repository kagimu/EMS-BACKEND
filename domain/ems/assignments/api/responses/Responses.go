package responses

import "time"

// UserAssignmentResponse is the API representation of an assignment.
//
// Privacy: only the referenced user's number is exposed — never their name,
// contact details, or any authentication data. No personal data is returned.
type UserAssignmentResponse struct {
	Number                 string     `json:"number"`
	OrganisationNumber     string     `json:"organisation_number"`
	UserNumber             string     `json:"user_number"`
	OrganisationUnitNumber string     `json:"organisation_unit_number"`
	AssignmentType         string     `json:"assignment_type"`
	PositionTitle          string     `json:"position_title"`
	ValidFrom              *time.Time `json:"valid_from"`
	ValidTo                *time.Time `json:"valid_to"`
	Status                 string     `json:"status"`
	AddedByNumber          string     `json:"added_by_number"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}
