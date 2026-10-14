package responses

import "time"

// DecisionAuthorityResponse is the API representation of an authority definition.
type DecisionAuthorityResponse struct {
	Number                 string     `json:"number"`
	OrganisationNumber     string     `json:"organisation_number"`
	DecisionType           string     `json:"decision_type"`
	ApprovalLevel          int        `json:"approval_level"`
	MinAmount              *float64   `json:"min_amount"`
	MaxAmount              *float64   `json:"max_amount"`
	RoleNumber             string     `json:"role_number"`
	StakeholderType        string     `json:"stakeholder_type"`
	OrganisationUnitNumber string     `json:"organisation_unit_number"`
	Status                 string     `json:"status"`
	ValidFrom              *time.Time `json:"valid_from"`
	ValidTo                *time.Time `json:"valid_to"`
	AddedByNumber          string     `json:"added_by_number"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

// AuthorityAnchorResponse is one anchor (role or stakeholder type) that satisfies
// an approval level.
//
// Foundation D (ratified D5): an anchor's identity is **type + value + unit**. The
// unit is the authority row's organisation-unit scope; an empty unit means
// organisation-wide. Two anchors are the same anchor only when all three parts
// are equal, so FINANCE_MANAGER@KAMPALA and FINANCE_MANAGER@JINJA never collapse.
type AuthorityAnchorResponse struct {
	Type                   string `json:"type"` // "role" or "stakeholder_type"
	Value                  string `json:"value"`
	OrganisationUnitNumber string `json:"organisation_unit_number"`
}

// AuthorityLevelResponse is one required approval level in the resolved path.
type AuthorityLevelResponse struct {
	ApprovalLevel int                       `json:"approval_level"`
	Anchors       []AuthorityAnchorResponse `json:"anchors"`
}

// AuthorityResolutionResponse is the read-only description of who is required to
// approve a decision for a given amount. It describes a path; it does not create
// or execute anything.
type AuthorityResolutionResponse struct {
	Configured               bool                     `json:"configured"`
	DecisionType             string                   `json:"decision_type"`
	Amount                   float64                  `json:"amount"`
	OrganisationUnitNumber   string                   `json:"organisation_unit_number"`
	ExceedsDeclaredAuthority bool                     `json:"exceeds_declared_authority"`
	RequiredPath             []AuthorityLevelResponse `json:"required_path"`
}
