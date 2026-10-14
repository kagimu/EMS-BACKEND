package models

import (
	"time"

	"gorm.io/gorm"
)

// Installation lifecycle statuses. Client self-deactivation intentionally does
// NOT change the status: a client-deactivated installation remains ACTIVE and
// therefore still consumes capacity until an explicit control-plane release
// (Foundation F §8, §24 detail 2).
const (
	InstallationStatusActive   = "ACTIVE"
	InstallationStatusReleased = "RELEASED"
	InstallationStatusRevoked  = "REVOKED"
)

// Installation is a licence-bound activation record: one row per activation of
// an EMS installation against a licence. It is organisation-bound and
// licence-bound by value (no foreign keys into existing tables).
//
// The installation identity is an opaque, client-generated identifier used for
// recognition only. It is never an authority grant, and it is never silently
// transferred between licences (F-Q3, §13).
type Installation struct {
	gorm.Model
	Number               string    `json:"number" gorm:"size:36;unique;not null"`
	InstallationIdentity string    `json:"installation_identity" gorm:"size:128;not null;uniqueIndex:idx_ems_installations_licence_identity"`
	OrganisationNumber   string    `json:"organisation_number" gorm:"size:36;not null;index:idx_ems_installations_org_status"`
	LicenceNumber        string    `json:"licence_number" gorm:"size:36;not null;uniqueIndex:idx_ems_installations_licence_identity;index:idx_ems_installations_licence_status"`
	ActivatedAt          time.Time `json:"activated_at" gorm:"not null"`
	// ActivatedBy records the client actor of the original activation (§17/§18).
	ActivatedBy *string    `json:"activated_by" gorm:"size:36"`
	LastSeenAt  *time.Time `json:"last_seen_at"`
	Status      string     `json:"status" gorm:"size:20;not null;index:idx_ems_installations_org_status;index:idx_ems_installations_licence_status"`
	// ClientDeactivatedAt/By mark a permissible client self-deactivation. The
	// status stays ACTIVE so the installation remains capacity-holding.
	ClientDeactivatedAt *time.Time `json:"client_deactivated_at"`
	ClientDeactivatedBy *string    `json:"client_deactivated_by" gorm:"size:36"`
	// Control-plane release is the ONLY action that frees capacity (F-Q3A).
	ReleasedBy    *string    `json:"released_by" gorm:"size:36"`
	ReleasedAt    *time.Time `json:"released_at"`
	ReleaseReason *string    `json:"release_reason" gorm:"size:128"`
	// ReplacedByNumber optionally links an owner-released installation to the
	// installation that replaced it, without altering the historical record.
	ReplacedByNumber *string `json:"replaced_by_number" gorm:"size:36"`
}

// TableName pins the table to ems_installations.
func (Installation) TableName() string {
	return "ems_installations"
}
