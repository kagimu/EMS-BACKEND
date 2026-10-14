package models

import (
	"time"

	"gorm.io/gorm"
)

// Stored lifecycle statuses. These are the ONLY persisted status values.
// Derived states (IN_GRACE, EXPIRED) are computed at read/request time from
// the validity window and are never stored (Foundation F, F-Q2).
const (
	StatusActive    = "ACTIVE"
	StatusSuspended = "SUSPENDED"
	StatusRevoked   = "REVOKED"
)

// Derived effective statuses (never persisted).
const (
	EffectiveStatusActive    = "ACTIVE"
	EffectiveStatusInGrace   = "IN_GRACE"
	EffectiveStatusExpired   = "EXPIRED"
	EffectiveStatusSuspended = "SUSPENDED"
	EffectiveStatusRevoked   = "REVOKED"
)

// DefaultGraceDays is the ratified 14-day grace period (F-Q2). It is stored
// explicitly on each licence so the derivation is deterministic.
const DefaultGraceDays = 14

// Licence is the Foundation F licence registry record: one row per issued
// licence. It is organisation-bound (the existing eos.Entity, referenced by
// value) and carries the opaque licence key, validity window, installation
// allowance and module entitlements.
//
// This is the only table introduced by Foundation F Increment 1 (ems_licences).
// There are no foreign keys into existing tables — references are by value.
type Licence struct {
	gorm.Model
	Number             string    `json:"number" gorm:"size:36;unique;not null"`
	LicenceKey         string    `json:"licence_key" gorm:"size:64;unique;not null"`
	OrganisationNumber string    `json:"organisation_number" gorm:"size:36;not null;index"`
	ValidFrom          time.Time `json:"valid_from" gorm:"not null"`
	ValidTo            time.Time `json:"valid_to" gorm:"not null"`
	GraceDays          int       `json:"grace_days" gorm:"not null;default:14"`
	MaxInstallations   int       `json:"max_installations" gorm:"not null"`
	LicensedAddons     string    `json:"licensed_addons" gorm:"size:500;not null"`
	Status             string    `json:"status" gorm:"size:20;not null"`
	// ActiveOrgNumber is the one-active-licence guard: it equals
	// OrganisationNumber while the stored status is ACTIVE and is NULL
	// otherwise. The unique index makes a second ACTIVE licence for the same
	// organisation impossible even under concurrent issuance (F-Q7).
	ActiveOrgNumber  *string `json:"active_org_number" gorm:"size:36;unique"`
	ReplacedByNumber *string `json:"replaced_by_number" gorm:"size:36"`
	IssuedBy         string  `json:"issued_by" gorm:"size:36;not null"`
	// LastActedBy records the most recent control-plane actor for a licence
	// transition (configure, renew, suspend, revoke, replace). Durable actor
	// facts are kept on F's own rows; Foundation G will reference them (§18).
	LastActedBy *string `json:"last_acted_by" gorm:"size:36"`
	// RenewalRequestedBy/At record the client renewal request (§18/§19). A
	// request never changes validity or any other licence state.
	RenewalRequestedBy *string    `json:"renewal_requested_by" gorm:"size:36"`
	RenewalRequestedAt *time.Time `json:"renewal_requested_at"`
}

// TableName pins the table to ems_licences.
func (Licence) TableName() string {
	return "ems_licences"
}
