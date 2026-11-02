package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Ticket statuses. Transitions are validated in the service.
const (
	StatusOpen               = "OPEN"
	StatusInProgress         = "IN_PROGRESS"
	StatusWaitingForCustomer = "WAITING_FOR_CUSTOMER"
	StatusWaitingForSupport  = "WAITING_FOR_SUPPORT"
	StatusResolved           = "RESOLVED"
	StatusClosed             = "CLOSED"
)

// Priority values.
const (
	PriorityLow      = "LOW"
	PriorityMedium   = "MEDIUM"
	PriorityHigh     = "HIGH"
	PriorityCritical = "CRITICAL"
)

// Message visibility.
const (
	VisibilityCustomer = "CUSTOMER"
	VisibilityInternal = "INTERNAL"
)

// Author roles.
const (
	AuthorOwner    = "OWNER"
	AuthorCustomer = "CUSTOMER"
	AuthorManager  = "MANAGER"
	AuthorAgent    = "AGENT"
)

// Support staff roles. These are support-workspace authorities, deliberately
// separate from ordinary client organisation roles.
const (
	SupportRoleManager    = "MANAGER"
	SupportRoleAgent      = "AGENT"
	SupportStatusActive   = "ACTIVE"
	SupportStatusInactive = "INACTIVE"
)

// Notification delivery states. CAPTURED means stored in-app; no external
// provider is configured in this stage.
const (
	DeliveryInApp         = "IN_APP"
	DeliveryQueued        = "QUEUED"
	DeliverySent          = "SENT"
	DeliveryFailed        = "FAILED"
	DeliveryNotConfigured = "EMAIL_NOT_CONFIGURED"
)

// Categories.
var Categories = []string{
	"ACCOUNT_ACCESS",
	"PASSWORD_RESET",
	"LICENCE_ACTIVATION",
	"LICENCE_EXPIRY",
	"MODULE_ACCESS",
	"INSTALLATION_PROBLEM",
	"USER_MANAGEMENT",
	"TECHNICAL_ERROR",
	"DATA_CONFIGURATION",
	"GENERAL_ENQUIRY",
	"OTHER",
}

var Priorities = []string{PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical}
var Statuses = []string{StatusOpen, StatusInProgress, StatusWaitingForCustomer, StatusWaitingForSupport, StatusResolved, StatusClosed}

// Ticket is an additive support record. It references organisations, users,
// licences and installations by value; it does not create duplicate customer or
// licensing entities.
type Ticket struct {
	gorm.Model
	Number              string     `json:"number" gorm:"size:36;unique;not null"`
	OrganisationNumber  string     `json:"organisation_number" gorm:"size:36;not null;index"`
	RequestedByNumber   string     `json:"requested_by_number" gorm:"size:36;index"`
	Subject             string     `json:"subject" gorm:"size:191;not null"`
	Description         string     `json:"description" gorm:"type:text"`
	Category            string     `json:"category" gorm:"size:40;not null;index"`
	Priority            string     `json:"priority" gorm:"size:20;not null;index"`
	Status              string     `json:"status" gorm:"size:30;not null;index"`
	AssignedToNumber    *string    `json:"assigned_to_number,omitempty" gorm:"size:36;index"`
	LicenceNumber       *string    `json:"licence_number,omitempty" gorm:"size:36;index"`
	InstallationNumber  *string    `json:"installation_number,omitempty" gorm:"size:36;index"`
	ResolvedAt          *time.Time `json:"resolved_at,omitempty"`
	ResolutionSummary   string     `json:"resolution_summary" gorm:"type:text"`
	FirstResponseAt     *time.Time `json:"first_response_at,omitempty"`
	SLAFirstResponseDue *time.Time `json:"sla_first_response_due,omitempty"`
	SLAResolutionDue    *time.Time `json:"sla_resolution_due,omitempty"`
	Escalated           bool       `json:"escalated" gorm:"not null;default:false;index"`
	EscalatedAt         *time.Time `json:"escalated_at,omitempty"`
	EscalatedBy         *string    `json:"escalated_by,omitempty" gorm:"size:36"`
	EscalationReason    string     `json:"escalation_reason" gorm:"type:text"`
}

func (Ticket) TableName() string { return "ems_support_tickets" }

// Message is a ticket conversation item. Internal messages are never returned
// by customer-facing APIs.
type Message struct {
	gorm.Model
	TicketNumber   string `json:"ticket_number" gorm:"size:36;not null;index"`
	AuthorNumber   string `json:"author_number" gorm:"size:36;not null"`
	AuthorRole     string `json:"author_role" gorm:"size:20;not null"`
	Visibility     string `json:"visibility" gorm:"size:20;not null;index"`
	Body           string `json:"body" gorm:"type:text;not null"`
	DeliveryStatus string `json:"delivery_status" gorm:"size:30;not null"`
}

func (Message) TableName() string { return "ems_support_ticket_messages" }

// Audit is append-only. No update or delete endpoint exists.
type Audit struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	TicketNumber string    `json:"ticket_number" gorm:"size:36;not null;index"`
	ActorNumber  string    `json:"actor_number" gorm:"size:36;not null"`
	Action       string    `json:"action" gorm:"size:60;not null;index"`
	BeforeState  string    `json:"before_state" gorm:"type:text"`
	AfterState   string    `json:"after_state" gorm:"type:text"`
	Details      string    `json:"details" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at"`
}

func (Audit) TableName() string { return "ems_support_ticket_audits" }

// Staff is a support-workspace authority mapped to a normal EMS user. It is not
// an identity system and grants no client organisation role.
type Staff struct {
	gorm.Model
	UserNumber    string `json:"user_number" gorm:"size:36;unique;not null"`
	Role          string `json:"role" gorm:"size:20;not null;index"`
	Status        string `json:"status" gorm:"size:20;not null;index"`
	AddedByNumber string `json:"added_by_number" gorm:"size:36"`
}

func (Staff) TableName() string { return "ems_support_staff" }

// Attachment metadata. Files are stored outside the public web root and served
// only through an authorised download endpoint.
type Attachment struct {
	gorm.Model
	Number           string `json:"number" gorm:"size:36;unique;not null"`
	TicketNumber     string `json:"ticket_number" gorm:"size:36;not null;index"`
	UploadedByNumber string `json:"uploaded_by_number" gorm:"size:36;not null"`
	OriginalName     string `json:"original_name" gorm:"size:255;not null"`
	StoredName       string `json:"-" gorm:"size:100;not null"`
	ContentType      string `json:"content_type" gorm:"size:100;not null"`
	SizeBytes        int64  `json:"size_bytes" gorm:"not null"`
	Visibility       string `json:"visibility" gorm:"size:20;not null;index"`
}

func (Attachment) TableName() string { return "ems_support_ticket_attachments" }

// Notification is a delivery record for support events. DedupeKey is unique so
// repeated triggers do not spam recipients.
type Notification struct {
	gorm.Model
	TicketNumber        string `json:"ticket_number" gorm:"size:36;not null;index"`
	RecipientUserNumber string `json:"recipient_user_number" gorm:"size:36;not null;index"`
	NotificationType    string `json:"notification_type" gorm:"size:60;not null"`
	DedupeKey           string `json:"dedupe_key" gorm:"size:191;unique;not null"`
	DeliveryStatus      string `json:"delivery_status" gorm:"size:40;not null"`
	Attempts            int    `json:"attempts" gorm:"not null;default:0"`
	LastError           string `json:"last_error" gorm:"type:text"`
}

func (Notification) TableName() string { return "ems_support_notifications" }

// SLATarget is the owner-configurable target per priority. Durations are
// calendar minutes in this stage; no business calendar or timezone is applied.
type SLATarget struct {
	gorm.Model
	Priority             string `json:"priority" gorm:"size:20;unique;not null"`
	FirstResponseMinutes int    `json:"first_response_minutes" gorm:"not null"`
	ResolutionMinutes    int    `json:"resolution_minutes" gorm:"not null"`
	UpdatedByNumber      string `json:"updated_by_number" gorm:"size:36"`
}

func (SLATarget) TableName() string { return "ems_support_sla_targets" }

func IsSupportRole(value string) bool {
	return value == SupportRoleManager || value == SupportRoleAgent
}

func IsCategory(value string) bool {
	for _, category := range Categories {
		if value == category {
			return true
		}
	}
	return false
}

func IsPriority(value string) bool {
	for _, priority := range Priorities {
		if value == priority {
			return true
		}
	}
	return false
}

func IsStatus(value string) bool {
	for _, status := range Statuses {
		if value == status {
			return true
		}
	}
	return false
}

// ValidateTransition enforces the support lifecycle. Reopening is allowed from
// RESOLVED only; CLOSED is terminal.
func ValidateTransition(from, to string) error {
	if from == to {
		return nil
	}
	allowed := map[string]map[string]bool{
		StatusOpen: {
			StatusInProgress: true, StatusWaitingForCustomer: true, StatusWaitingForSupport: true,
			StatusResolved: true, StatusClosed: true,
		},
		StatusInProgress: {
			StatusWaitingForCustomer: true, StatusWaitingForSupport: true, StatusResolved: true, StatusClosed: true,
		},
		StatusWaitingForCustomer: {
			StatusInProgress: true, StatusWaitingForSupport: true, StatusResolved: true, StatusClosed: true,
		},
		StatusWaitingForSupport: {
			StatusInProgress: true, StatusWaitingForCustomer: true, StatusResolved: true, StatusClosed: true,
		},
		StatusResolved: {StatusClosed: true, StatusOpen: true, StatusInProgress: true},
		StatusClosed:   {},
	}
	if !allowed[from][to] {
		return errors.New("invalid ticket status transition")
	}
	return nil
}

func NormaliseCategory(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if !IsCategory(value) {
		return "", errors.New("invalid ticket category")
	}
	return value, nil
}

func NormalisePriority(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if !IsPriority(value) {
		return "", errors.New("invalid ticket priority")
	}
	return value, nil
}

func NormaliseStatus(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if !IsStatus(value) {
		return "", errors.New("invalid ticket status")
	}
	return value, nil
}
