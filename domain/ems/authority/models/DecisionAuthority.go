package models

import (
	"time"

	"gorm.io/gorm"
)

// DecisionAuthority declares who may make a particular organisational decision:
// for one decision type, at an ordered approval level, optionally scoped to an
// organisation unit, up to an amount limit, anchored on a system role OR a
// stakeholder type.
//
// Foundation D: this is the only table introduced by this release
// (ems_decision_authorities). It DEFINES authority and is used by the read-only
// resolver. It does not create approvals and does not execute workflow
// (Foundation E).
type DecisionAuthority struct {
	gorm.Model
	Number                 string     `json:"number" gorm:"size:36;unique;not null"`
	OrganisationNumber     string     `json:"organisation_number" gorm:"size:36;not null;index"`
	DecisionType           string     `json:"decision_type" gorm:"size:50;not null;index"`
	ApprovalLevel          int        `json:"approval_level" gorm:"not null;index"`
	MinAmount              *float64   `json:"min_amount" gorm:"type:decimal(18,2)"`
	MaxAmount              *float64   `json:"max_amount" gorm:"type:decimal(18,2)"`
	RoleNumber             string     `json:"role_number" gorm:"size:36;index"`
	StakeholderType        string     `json:"stakeholder_type" gorm:"size:30;index"`
	OrganisationUnitNumber string     `json:"organisation_unit_number" gorm:"size:36;index"`
	Status                 string     `json:"status" gorm:"size:20;not null"`
	ValidFrom              *time.Time `json:"valid_from"`
	ValidTo                *time.Time `json:"valid_to"`
	AddedByNumber          string     `json:"added_by_number" gorm:"size:36;not null"`
}

// TableName pins the table to ems_decision_authorities.
func (DecisionAuthority) TableName() string {
	return "ems_decision_authorities"
}
