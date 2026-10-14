package responses

import "time"

// ApprovalAnchorResponse is one snapshotted authority anchor on a step.
//
// Foundation D (ratified D5): the unit is part of the anchor's identity, so a
// unit-scoped anchor is distinguishable from an organisation-wide one (empty unit).
type ApprovalAnchorResponse struct {
	Type                   string `json:"type"` // "role" or "stakeholder_type"
	Value                  string `json:"value"`
	OrganisationUnitNumber string `json:"organisation_unit_number"`
}

// ApprovalAssignmentResponse is one eligible candidate (and its decision, if any).
type ApprovalAssignmentResponse struct {
	Number     string     `json:"number"`
	UserNumber string     `json:"user_number"`
	Status     string     `json:"status"`
	Comment    string     `json:"comment"`
	DecidedAt  *time.Time `json:"decided_at"`
}

// ApprovalStepResponse is one ordered authority level.
//
// Actionable is true only when the request is PENDING and this step is the
// request's current level and the step itself is PENDING.
type ApprovalStepResponse struct {
	Number        string                       `json:"number"`
	ApprovalLevel int                          `json:"approval_level"`
	Status        string                       `json:"status"`
	Actionable    bool                         `json:"actionable"`
	Anchors       []ApprovalAnchorResponse     `json:"anchors"`
	Assignments   []ApprovalAssignmentResponse `json:"assignments"`
}

// ApprovalRequestResponse is the API representation of an approval request.
type ApprovalRequestResponse struct {
	Number                 string   `json:"number"`
	OrganisationNumber     string   `json:"organisation_number"`
	OrganisationUnitNumber string   `json:"organisation_unit_number"`
	Module                 string   `json:"module"`
	ResourceType           string   `json:"resource_type"`
	ResourceNumber         string   `json:"resource_number"`
	DecisionType           string   `json:"decision_type"`
	Amount                 *float64 `json:"amount"`
	Status                 string   `json:"status"`
	CurrentLevel           int      `json:"current_level"`
	// Stalled is Foundation E3 (ratified D2/D3). It is DERIVED AT READ TIME and is
	// never stored: it is true when the request is PENDING, its current level still
	// awaits a decision, and NONE of that level's snapshotted candidates is
	// currently eligible to act. It is informational only — it never substitutes,
	// reassigns or re-resolves anything, and the action-time guard remains the sole
	// authority for acting.
	Stalled bool `json:"stalled"`
	// DeadlineAt is the Foundation E4 (ratified E4-D4) EFFECTIVE advisory
	// deadline for this request: the per-request override when present, otherwise
	// the configured default applied to CreatedAt. Always UTC (E4-D5). Advisory
	// only — it never changes the request status.
	DeadlineAt *time.Time `json:"deadline_at,omitempty"`
	// Overdue is Foundation E4 (ratified E4-D1/D2/D6, C1 Option B, C4). It is
	// DERIVED AT READ TIME and never stored: true only when the request is PENDING
	// and the effective deadline has passed (nowUTC > DeadlineAt). It is
	// informational only — it never changes status, authority or eligibility, and
	// it never grants the right to act.
	Overdue bool `json:"overdue"`
	// Reminder is Foundation E4 (ratified C1 Option B, C4) — the "approaching
	// deadline" condition. It is DERIVED AT READ TIME and never stored: true only
	// when the request is PENDING, the effective deadline is still in the future
	// (nowUTC <= DeadlineAt) and no more than 24 hours away. Reminder and Overdue
	// are mutually exclusive; a non-PENDING request is neither. Informational
	// only — it never changes status, authority or eligibility.
	Reminder              bool                   `json:"reminder"`
	RequestedByNumber     string                 `json:"requested_by_number"`
	PreviousRequestNumber string                 `json:"previous_request_number"`
	Steps                 []ApprovalStepResponse `json:"steps"`
	CreatedAt             time.Time              `json:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
}

// ApprovalTaskResponse is one actionable assignment for the caller's inbox.
type ApprovalTaskResponse struct {
	RequestNumber          string   `json:"request_number"`
	StepNumber             string   `json:"step_number"`
	AssignmentNumber       string   `json:"assignment_number"`
	ApprovalLevel          int      `json:"approval_level"`
	Module                 string   `json:"module"`
	ResourceType           string   `json:"resource_type"`
	ResourceNumber         string   `json:"resource_number"`
	DecisionType           string   `json:"decision_type"`
	Amount                 *float64 `json:"amount"`
	OrganisationUnitNumber string   `json:"organisation_unit_number"`
	RequestStatus          string   `json:"request_status"`
	// DeadlineAt / Reminder / Overdue are the Foundation E4 (ratified C1 Option B,
	// C4) READ-TIME derived deadline conditions for this task's request, surfaced
	// so the inbox can distinguish normal / approaching / overdue without any
	// durable notification. They are never stored and create no notification.
	DeadlineAt *time.Time `json:"deadline_at,omitempty"`
	Reminder   bool       `json:"reminder"`
	Overdue    bool       `json:"overdue"`
}

// DeadlineConfigResponse is the Foundation E4 (ratified E4-D4) workflow deadline
// configuration for one organisation. It is advisory configuration only and
// never changes an existing request's established deadline.
type DeadlineConfigResponse struct {
	OrganisationNumber   string `json:"organisation_number"`
	DefaultDeadlineHours int    `json:"default_deadline_hours"`
}
