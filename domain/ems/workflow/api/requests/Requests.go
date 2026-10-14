package requests

import "time"

// CreateApprovalRequest creates a generic approval request.
//
// The organisation is NEVER supplied by the client; it comes from the
// authenticated active-organisation context. The amount is optional and is used
// only as input to Foundation D's authority resolution.
type CreateApprovalRequest struct {
	Module                 string   `json:"module"`
	ResourceType           string   `json:"resource_type"`
	ResourceNumber         string   `json:"resource_number"`
	DecisionType           string   `json:"decision_type"`
	Amount                 *float64 `json:"amount"`
	OrganisationUnitNumber string   `json:"organisation_unit_number"`

	// PreviousRequestNumber optionally links this creation to the RETURNED
	// request it corrects (Foundation E2 resubmission). It is validated
	// server-side and is informational only: it never affects authority,
	// eligibility, permissions or organisation.
	PreviousRequestNumber string `json:"previous_request_number"`

	// DeadlineAt is the optional Foundation E4 (ratified E4-D4) per-request
	// advisory deadline override. When omitted, the configured default applies.
	// It is validated server-side (must be in the future, UTC) and is advisory
	// only: it never affects authority, eligibility, permissions or status.
	DeadlineAt *time.Time `json:"deadline_at"`
}

// ApprovalDecisionRequest is the payload for approve / reject.
type ApprovalDecisionRequest struct {
	Comment string `json:"comment"`
}

// UpdateDeadlineConfigRequest is the Foundation E4 (ratified E4-D4) payload for
// setting the organisation's default approval deadline, in hours. The
// organisation is never supplied by the client; it comes from the authenticated
// active-organisation context.
type UpdateDeadlineConfigRequest struct {
	DefaultDeadlineHours int `json:"default_deadline_hours"`
}
