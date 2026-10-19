package responses

import "time"

// ReportingLineResponse is the API representation of a reporting relationship.
type ReportingLineResponse struct {
	Number                      string     `json:"number"`
	OrganisationNumber          string     `json:"organisation_number"`
	SubordinateAssignmentNumber string     `json:"subordinate_assignment_number"`
	ManagerAssignmentNumber     string     `json:"manager_assignment_number"`
	ValidFrom                   *time.Time `json:"valid_from"`
	ValidTo                     *time.Time `json:"valid_to"`
	Status                      string     `json:"status"`
	Effective                   bool       `json:"effective"`
	EffectiveState              string     `json:"effective_state"`
	AddedByNumber               string     `json:"added_by_number"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`
}
