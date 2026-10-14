package models

import (
	"time"

	"gorm.io/gorm"
)

// Approval request statuses (request-level state machine).
const (
	RequestStatusPending   = "PENDING"
	RequestStatusApproved  = "APPROVED"
	RequestStatusRejected  = "REJECTED"
	RequestStatusCancelled = "CANCELLED"
	RequestStatusBlocked   = "BLOCKED"
	// RequestStatusReturned is Foundation E2 (ratified D4): an eligible approver
	// of the current step sent the request back to the requester for correction.
	// It is terminal for this request; the requester's remedy is resubmission as
	// a NEW request (ratified D5-D7).
	RequestStatusReturned = "RETURNED"
)

// Approval step statuses (one ordered authority level).
//
// Foundation E1 is strictly sequential: only the request's current level is
// actionable. Steps that have not been reached are PENDING but are NOT
// actionable (guarded by request.current_level).
const (
	StepStatusPending   = "PENDING"
	StepStatusApproved  = "APPROVED"
	StepStatusRejected  = "REJECTED"
	StepStatusCancelled = "CANCELLED"
	StepStatusBlocked   = "BLOCKED"
	// StepStatusReturned mirrors RequestStatusReturned for the step that was
	// returned (Foundation E2).
	StepStatusReturned = "RETURNED"
)

// Approval assignment statuses. An assignment row records one eligible
// candidate for a step and carries its write-once decision (BD-E1: ANY — the
// first accepted decision completes the step).
const (
	AssignmentStatusUndecided = "UNDECIDED"
	AssignmentStatusApproved  = "APPROVED"
	AssignmentStatusRejected  = "REJECTED"
	// AssignmentStatusReturned is the returning approver's write-once decision
	// row, carrying the mandatory reason in Comment (Foundation E2).
	AssignmentStatusReturned = "RETURNED"
)

// ApprovalRequest is a generic business approval awaiting a decision.
//
// Foundation E1: one of the three additive tables (ems_approval_requests). It
// references the business object generically (module / resource_type /
// resource_number, by value) and never dereferences it. The organisation is
// always taken from the authenticated context.
type ApprovalRequest struct {
	gorm.Model
	Number                 string   `json:"number" gorm:"size:36;unique;not null"`
	OrganisationNumber     string   `json:"organisation_number" gorm:"size:36;not null;index"`
	OrganisationUnitNumber string   `json:"organisation_unit_number" gorm:"size:36;index"`
	Module                 string   `json:"module" gorm:"size:50;not null;index"`
	ResourceType           string   `json:"resource_type" gorm:"size:50;not null;index"`
	ResourceNumber         string   `json:"resource_number" gorm:"size:36;not null;index"`
	DecisionType           string   `json:"decision_type" gorm:"size:50;not null;index"`
	Amount                 *float64 `json:"amount" gorm:"type:decimal(18,2)"`
	Status                 string   `json:"status" gorm:"size:20;not null;index"`
	CurrentLevel           int      `json:"current_level" gorm:"not null;default:0"`
	RequestedByNumber      string   `json:"requested_by_number" gorm:"size:36;not null;index"`
	AddedByNumber          string   `json:"added_by_number" gorm:"size:36;not null"`
	// PreviousRequestNumber links a ratified-E2 resubmission to the request that
	// preceded it (the RETURNED request it corrects). It is INFORMATIONAL ONLY
	// and must never be consulted for authority, eligibility, permissions,
	// organisation or actionability. Empty for an ordinary creation.
	PreviousRequestNumber string `json:"previous_request_number" gorm:"size:36;index"`

	// DeadlineAt is the Foundation E4 (ratified E4-D4) optional per-request
	// advisory deadline override. It is NULLable: when empty the effective
	// deadline is the configured default applied to CreatedAt. It is ADVISORY
	// ONLY (ratified E4-D2) — it never changes the request status, authority,
	// eligibility or the snapshot, and `overdue` is derived at read time and
	// never stored (ratified E4-D6). Stored and compared in UTC (E4-D5).
	DeadlineAt *time.Time `json:"deadline_at" gorm:"index"`
}

func (ApprovalRequest) TableName() string { return "ems_approval_requests" }

// ApprovalStep is one ordered authority level of a request, snapshotted when the
// request is created (BD-E3/BD-E14: D is never re-run for an existing request).
//
// Anchors are stored as a compact serialised list because Foundation D returns a
// LIST of anchors per level (e.g. a role and a stakeholder type at the same
// level); the eligible-candidate set is the union across those anchors.
// Format: "type:value|type:value" (e.g. "role:<uuid>|stakeholder_type:owner").
type ApprovalStep struct {
	gorm.Model
	Number             string `json:"number" gorm:"size:36;unique;not null"`
	RequestNumber      string `json:"request_number" gorm:"size:36;not null;index;uniqueIndex:idx_ems_approval_steps_request_level"`
	OrganisationNumber string `json:"organisation_number" gorm:"size:36;not null;index"`
	ApprovalLevel      int    `json:"approval_level" gorm:"not null;index;uniqueIndex:idx_ems_approval_steps_request_level"`
	Anchors            string `json:"anchors" gorm:"size:500;not null"`
	Status             string `json:"status" gorm:"size:20;not null;index"`
}

func (ApprovalStep) TableName() string { return "ems_approval_steps" }

// ApprovalAssignment is one eligible candidate for a step.
//
// Under ANY semantics (BD-E1) any one candidate may act; the first accepted
// decision completes the step and the remaining candidate rows stay UNDECIDED
// forever. A decided row is never updated again (write-once decision record).
type ApprovalAssignment struct {
	gorm.Model
	Number             string     `json:"number" gorm:"size:36;unique;not null"`
	StepNumber         string     `json:"step_number" gorm:"size:36;not null;index;uniqueIndex:idx_ems_approval_assignments_step_user"`
	RequestNumber      string     `json:"request_number" gorm:"size:36;not null;index"`
	OrganisationNumber string     `json:"organisation_number" gorm:"size:36;not null;index"`
	UserNumber         string     `json:"user_number" gorm:"size:36;not null;index;uniqueIndex:idx_ems_approval_assignments_step_user"`
	Status             string     `json:"status" gorm:"size:20;not null;index"`
	Comment            string     `json:"comment" gorm:"size:1000"`
	DecidedAt          *time.Time `json:"decided_at"`
}

func (ApprovalAssignment) TableName() string { return "ems_approval_assignments" }

// ApprovalStepAnchor is one snapshotted authority anchor of an approval step.
//
// Foundation D (ratified D5 + D6): an anchor's identity is **type + value + unit**,
// and anchors are stored normalised — one row per anchor — rather than in a single
// serialised string. Two same-role anchors with different units therefore exist as
// two rows and can never collapse into one.
//
// OrganisationUnitNumber is the authority's scope for this anchor: empty means
// organisation-wide, in which case candidate resolution applies NO placement filter
// (ratified D7).
//
// This is the authoritative anchor record for steps created after Foundation D.
// Steps created before it have no rows here and continue to be read from
// ems_approval_steps.anchors, which is unit-less and therefore organisation-wide —
// so a legacy request is never retroactively narrowed (ratified D9 + D11).
//
// Value references only: no foreign keys, consistent with every other EMS table.
type ApprovalStepAnchor struct {
	gorm.Model
	Number                 string `json:"number" gorm:"size:36;unique;not null"`
	StepNumber             string `json:"step_number" gorm:"size:36;not null;index"`
	OrganisationNumber     string `json:"organisation_number" gorm:"size:36;not null;index"`
	AnchorType             string `json:"anchor_type" gorm:"size:20;not null"`
	AnchorValue            string `json:"anchor_value" gorm:"size:36;not null"`
	OrganisationUnitNumber string `json:"organisation_unit_number" gorm:"size:36"`
}

func (ApprovalStepAnchor) TableName() string { return "ems_approval_step_anchors" }
