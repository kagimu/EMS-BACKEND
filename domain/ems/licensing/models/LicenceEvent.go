package models

import (
	"time"

	"gorm.io/gorm"
)

// LicenceEventType is the closed vocabulary of Foundation F licensing facts.
// These are append-only durable facts recorded by Foundation F for Foundation
// G to reference later (§17/§18). Foundation F does not deliver notifications
// and does not run a second audit engine; it only records facts.
type LicenceEventType string

const (
	// Licence lifecycle facts.
	LicenceEventIssued     LicenceEventType = "LICENCE_ISSUED"
	LicenceEventConfigured LicenceEventType = "LICENCE_CONFIGURED"
	LicenceEventRenewed    LicenceEventType = "LICENCE_RENEWED"
	LicenceEventSuspended  LicenceEventType = "LICENCE_SUSPENDED"
	LicenceEventRevoked    LicenceEventType = "LICENCE_REVOKED"
	LicenceEventReplaced   LicenceEventType = "LICENCE_REPLACED"

	// Installation facts.
	InstallationEventActivated   LicenceEventType = "INSTALLATION_ACTIVATED"
	InstallationEventReactivated LicenceEventType = "INSTALLATION_REACTIVATED"
	InstallationEventDeactivated LicenceEventType = "INSTALLATION_DEACTIVATED"
	InstallationEventReleased    LicenceEventType = "INSTALLATION_RELEASED"

	// Client action facts.
	ClientEventRenewalRequested LicenceEventType = "CLIENT_RENEWAL_REQUESTED"
)

// LicenceEvent is the Foundation F append-only licensing fact record. Rows are
// never updated or deleted. It captures actor, organisation, action, timestamp
// and the related licence/installation so Foundation G can reference the fact
// later via entity_type/entity_id without duplicating workflow history.
//
// It is deliberately a thin fact log, not an audit engine: there is no
// processing, retention policy, delivery or consumer inside Foundation F.
type LicenceEvent struct {
	gorm.Model
	EventType          LicenceEventType `json:"event_type" gorm:"size:40;not null;index:idx_licence_events_type"`
	LicenceNumber      string           `json:"licence_number" gorm:"size:36;not null;index:idx_licence_events_licence"`
	OrganisationNumber string           `json:"organisation_number" gorm:"size:36;not null;index:idx_licence_events_org"`
	ActorUserNumber    string           `json:"actor_user_number" gorm:"size:36;not null;index:idx_licence_events_actor"`
	// Installation context (nullable for licence-level facts).
	InstallationNumber   *string `json:"installation_number,omitempty" gorm:"size:36;index:idx_licence_events_installation"`
	InstallationIdentity *string `json:"installation_identity,omitempty" gorm:"size:128"`
	// Before/after state snapshots as JSON text (nullable).
	BeforeState string `json:"before_state" gorm:"type:text"`
	AfterState  string `json:"after_state" gorm:"type:text"`
	// Reason where the action requires one (e.g. release reason).
	Reason *string `json:"reason,omitempty" gorm:"size:500"`
	// Client-initiated facts may carry client context.
	ClientInitiated bool   `json:"client_initiated" gorm:"not null;default:false"`
	ClientIP        string `json:"client_ip,omitempty" gorm:"size:45"`
	UserAgent       string `json:"user_agent,omitempty" gorm:"size:255"`
}

// TableName pins the table to ems_licence_events.
func (LicenceEvent) TableName() string {
	return "ems_licence_events"
}

// LicenceEventFilter narrows an event query. Nil fields are ignored.
type LicenceEventFilter struct {
	EventType          *LicenceEventType
	LicenceNumber      *string
	OrganisationNumber *string
	ActorUserNumber    *string
	InstallationNumber *string
	Since              *time.Time
	Until              *time.Time
	Limit              int
	Offset             int
}
