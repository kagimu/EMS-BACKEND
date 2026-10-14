package requests

import "time"

// CreateDecisionAuthorityRequest is the create payload.
//
// The organisation is never supplied by the client; it comes from the
// authenticated active-organisation context.
type CreateDecisionAuthorityRequest struct {
	DecisionType           string     `json:"decision_type"`
	ApprovalLevel          int        `json:"approval_level"`
	MinAmount              *float64   `json:"min_amount"`
	MaxAmount              *float64   `json:"max_amount"`
	RoleNumber             string     `json:"role_number"`
	StakeholderType        string     `json:"stakeholder_type"`
	OrganisationUnitNumber string     `json:"organisation_unit_number"`
	ValidFrom              *time.Time `json:"valid_from"`
	ValidTo                *time.Time `json:"valid_to"`
}

// UpdateDecisionAuthorityRequest replaces the editable fields of an authority row.
type UpdateDecisionAuthorityRequest struct {
	DecisionType           string     `json:"decision_type"`
	ApprovalLevel          int        `json:"approval_level"`
	MinAmount              *float64   `json:"min_amount"`
	MaxAmount              *float64   `json:"max_amount"`
	RoleNumber             string     `json:"role_number"`
	StakeholderType        string     `json:"stakeholder_type"`
	OrganisationUnitNumber string     `json:"organisation_unit_number"`
	ValidFrom              *time.Time `json:"valid_from"`
	ValidTo                *time.Time `json:"valid_to"`
}

// UpdateDecisionAuthorityStatusRequest is the status payload.
type UpdateDecisionAuthorityStatusRequest struct {
	Status string `json:"status"`
}
