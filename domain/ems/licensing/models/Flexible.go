package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	ScopeComputer     = "COMPUTER"
	ScopeInstallation = "INSTALLATION"
	ScopeOrganisation = "ORGANISATION"
	ScopeSite         = "SITE"
	ScopeUserSeat     = "USER_SEAT"
	ScopeCustom       = "CUSTOM"
)

const (
	FlexStatusActive   = "ACTIVE"
	FlexStatusInactive = "INACTIVE"
	UnitStatusActive   = "ACTIVE"
	UnitStatusReleased = "RELEASED"
	UnitStatusRevoked  = "REVOKED"
)

type LicenceScope struct {
	gorm.Model
	Number        string    `json:"number" gorm:"size:36;unique;not null"`
	LicenceNumber string    `json:"licence_number" gorm:"size:36;not null;uniqueIndex:idx_scope_unique"`
	ScopeType     string    `json:"scope_type" gorm:"size:30;not null;uniqueIndex:idx_scope_unique"`
	ScopeKey      string    `json:"scope_key" gorm:"size:100;not null;default:'';uniqueIndex:idx_scope_unique"`
	Quantity      int       `json:"quantity" gorm:"not null"`
	Status        string    `json:"status" gorm:"size:20;not null;index"`
	ValidFrom     time.Time `json:"valid_from"`
	ValidTo       time.Time `json:"valid_to"`
	ConfigJSON    string    `json:"config_json" gorm:"type:text"`
}

func (LicenceScope) TableName() string { return "ems_licence_scopes" }

type LicenceEntitlement struct {
	gorm.Model
	Number        string    `json:"number" gorm:"size:36;unique;not null"`
	LicenceNumber string    `json:"licence_number" gorm:"size:36;not null;uniqueIndex:idx_entitlement_unique"`
	Code          string    `json:"code" gorm:"size:60;not null;uniqueIndex:idx_entitlement_unique"`
	Quantity      int       `json:"quantity" gorm:"not null;default:1"`
	Status        string    `json:"status" gorm:"size:20;not null;index"`
	ValidFrom     time.Time `json:"valid_from"`
	ValidTo       time.Time `json:"valid_to"`
	ConfigJSON    string    `json:"config_json" gorm:"type:text"`
}

func (LicenceEntitlement) TableName() string { return "ems_licence_entitlements" }

type LicensedUnit struct {
	gorm.Model
	Number             string     `json:"number" gorm:"size:36;unique;not null"`
	LicenceNumber      string     `json:"licence_number" gorm:"size:36;not null;uniqueIndex:idx_unit_unique"`
	ScopeNumber        string     `json:"scope_number" gorm:"size:36;not null;index"`
	UnitType           string     `json:"unit_type" gorm:"size:30;not null;uniqueIndex:idx_unit_unique"`
	UnitIdentifier     string     `json:"unit_identifier" gorm:"size:150;not null;uniqueIndex:idx_unit_unique"`
	OrganisationNumber string     `json:"organisation_number" gorm:"size:36;not null;index"`
	SiteNumber         string     `json:"site_number" gorm:"size:36;index"`
	Status             string     `json:"status" gorm:"size:20;not null;index"`
	RegisteredAt       time.Time  `json:"registered_at"`
	RegisteredBy       string     `json:"registered_by" gorm:"size:36"`
	ReleasedAt         *time.Time `json:"released_at,omitempty"`
	ReleasedBy         *string    `json:"released_by,omitempty" gorm:"size:36"`
	ReleaseReason      string     `json:"release_reason" gorm:"type:text"`
	ConfigJSON         string     `json:"config_json" gorm:"type:text"`
}

func (LicensedUnit) TableName() string { return "ems_licensed_units" }

func IsScopeType(value string) bool {
	switch value {
	case ScopeComputer, ScopeInstallation, ScopeOrganisation, ScopeSite, ScopeUserSeat, ScopeCustom:
		return true
	}
	return false
}
