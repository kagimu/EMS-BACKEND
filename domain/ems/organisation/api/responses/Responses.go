package responses

import "time"

// OrganisationUnitResponse is the API representation of an EMS organisation unit.
type OrganisationUnitResponse struct {
	Number             string    `json:"number"`
	OrganisationNumber string    `json:"organisation_number"`
	ParentNumber       string    `json:"parent_number"`
	UnitType           string    `json:"unit_type"`
	Name               string    `json:"name"`
	Code               string    `json:"code"`
	Status             string    `json:"status"`
	AddedByNumber      string    `json:"added_by_number"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
