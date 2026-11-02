package responses

import (
	"time"

	licensingResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/api/responses"
)

type TicketSummary struct {
	Number               string     `json:"number"`
	OrganisationNumber   string     `json:"organisation_number"`
	OrganisationName     string     `json:"organisation_name,omitempty"`
	RequestedByNumber    string     `json:"requested_by_number,omitempty"`
	Subject              string     `json:"subject"`
	Category             string     `json:"category"`
	Priority             string     `json:"priority"`
	Status               string     `json:"status"`
	AssignedToNumber     string     `json:"assigned_to_number,omitempty"`
	LicenceNumber        string     `json:"licence_number,omitempty"`
	InstallationNumber   string     `json:"installation_number,omitempty"`
	ResolvedAt           *time.Time `json:"resolved_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	FirstResponseAt      *time.Time `json:"first_response_at,omitempty"`
	SLAFirstResponseDue  *time.Time `json:"sla_first_response_due,omitempty"`
	SLAResolutionDue     *time.Time `json:"sla_resolution_due,omitempty"`
	FirstResponseOverdue bool       `json:"first_response_overdue"`
	ResolutionOverdue    bool       `json:"resolution_overdue"`
	Escalated            bool       `json:"escalated"`
	EscalatedAt          *time.Time `json:"escalated_at,omitempty"`
	EscalationReason     string     `json:"escalation_reason,omitempty"`
}

type MessageResponse struct {
	ID             uint      `json:"id"`
	TicketNumber   string    `json:"ticket_number"`
	AuthorNumber   string    `json:"author_number"`
	AuthorRole     string    `json:"author_role"`
	Visibility     string    `json:"visibility"`
	Body           string    `json:"body"`
	DeliveryStatus string    `json:"delivery_status"`
	CreatedAt      time.Time `json:"created_at"`
}

type AuditResponse struct {
	ID           uint      `json:"id"`
	TicketNumber string    `json:"ticket_number"`
	ActorNumber  string    `json:"actor_number"`
	Action       string    `json:"action"`
	BeforeState  string    `json:"before_state,omitempty"`
	AfterState   string    `json:"after_state,omitempty"`
	Details      string    `json:"details,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// TechnicalContext is read-only and never contains licence keys or secrets.
type TechnicalContext struct {
	LicenceStatus           string     `json:"licence_status,omitempty"`
	LicenceValidFrom        *time.Time `json:"licence_valid_from,omitempty"`
	LicenceValidTo          *time.Time `json:"licence_valid_to,omitempty"`
	LicenceAddons           []string   `json:"licence_addons,omitempty"`
	InstallationStatus      string     `json:"installation_status,omitempty"`
	InstallationActivatedAt *time.Time `json:"installation_activated_at,omitempty"`
	InstallationLastSeenAt  *time.Time `json:"installation_last_seen_at,omitempty"`
}

type TicketDetail struct {
	Ticket      TicketSummary                             `json:"ticket"`
	Description string                                    `json:"description"`
	Resolution  string                                    `json:"resolution_summary,omitempty"`
	Messages    []MessageResponse                         `json:"messages"`
	Audit       []AuditResponse                           `json:"audit,omitempty"`
	Technical   TechnicalContext                          `json:"technical"`
	Events      []licensingResponses.LicenceEventResponse `json:"events,omitempty"`
	Attachments []AttachmentResponse                      `json:"attachments"`
}

type Overview struct {
	Total              int64 `json:"total"`
	Open               int64 `json:"open"`
	InProgress         int64 `json:"in_progress"`
	WaitingForCustomer int64 `json:"waiting_for_customer"`
	WaitingForSupport  int64 `json:"waiting_for_support"`
	HighPriority       int64 `json:"high_priority"`
	Critical           int64 `json:"critical"`
	Resolved           int64 `json:"resolved"`
	Closed             int64 `json:"closed"`
}

type Page[T any] struct {
	Items   []T   `json:"items"`
	Total   int64 `json:"total"`
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
}
