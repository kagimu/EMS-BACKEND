package requests

// CreateOrganisationUnitRequest is the create payload.
//
// The organisation is NEVER supplied by the client; it is taken from the
// authenticated active-organisation context.
type CreateOrganisationUnitRequest struct {
	Name         string `json:"name"`
	UnitType     string `json:"unit_type"`
	Code         string `json:"code"`
	ParentNumber string `json:"parent_number"`
}

// UpdateOrganisationUnitRequest is the update payload.
type UpdateOrganisationUnitRequest struct {
	Name         string `json:"name"`
	UnitType     string `json:"unit_type"`
	Code         string `json:"code"`
	ParentNumber string `json:"parent_number"`
}

// UpdateOrganisationUnitStatusRequest is the status payload.
type UpdateOrganisationUnitStatusRequest struct {
	Status string `json:"status"`
}
