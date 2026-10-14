package requests

import "time"

// CreateUserAssignmentRequest is the create payload.
//
// The organisation is never supplied by the client; it comes from the
// authenticated active-organisation context. There is no field for creating a
// user or a unit — both must already exist.
type CreateUserAssignmentRequest struct {
	UserNumber             string     `json:"user_number"`
	OrganisationUnitNumber string     `json:"organisation_unit_number"`
	AssignmentType         string     `json:"assignment_type"`
	PositionTitle          string     `json:"position_title"`
	ValidFrom              *time.Time `json:"valid_from"`
	ValidTo                *time.Time `json:"valid_to"`
}

// UpdateUserAssignmentRequest updates placement and lifecycle fields. The
// person is immutable: re-assigning a different person means ending this
// assignment and creating another.
type UpdateUserAssignmentRequest struct {
	OrganisationUnitNumber string     `json:"organisation_unit_number"`
	AssignmentType         string     `json:"assignment_type"`
	PositionTitle          string     `json:"position_title"`
	ValidFrom              *time.Time `json:"valid_from"`
	ValidTo                *time.Time `json:"valid_to"`
}

// UpdateUserAssignmentStatusRequest is the status payload.
type UpdateUserAssignmentStatusRequest struct {
	Status string `json:"status"`
}
