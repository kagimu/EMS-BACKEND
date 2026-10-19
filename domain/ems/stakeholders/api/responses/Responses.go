package responses

import "time"

// StakeholderResponse is the API representation of a stakeholder relationship.
//
// Privacy: the referenced user's number and display name (first/second/last
// name) are exposed for organisation-scoped presentation. Contact details,
// passwords and any authentication data are never exposed.
//
// CanAuthenticate is a derived, read-only flag (Foundation C2.3): it is true
// only when the referenced user currently has a usable login (an Email contact
// AND a password row AND an account status of "Active"). It exposes no contact
// or credential information and does not alter authentication behaviour.
type StakeholderResponse struct {
	Number             string     `json:"number"`
	OrganisationNumber string     `json:"organisation_number"`
	UserNumber         string     `json:"user_number"`
	DisplayName        string     `json:"display_name"`
	StakeholderType    string     `json:"stakeholder_type"`
	Title              string     `json:"title"`
	SharePercentage    *float64   `json:"share_percentage"`
	Shares             *int64     `json:"shares"`
	ValidFrom          *time.Time `json:"valid_from"`
	ValidTo            *time.Time `json:"valid_to"`
	Status             string     `json:"status"`
	CanAuthenticate    bool       `json:"can_authenticate"`
	AddedByNumber      string     `json:"added_by_number"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
