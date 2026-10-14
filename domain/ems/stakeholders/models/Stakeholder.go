package models

import (
	"time"

	"gorm.io/gorm"
)

// Stakeholder types. This is an open set: the list below is the initial
// catalogue and can be extended without changing the model. It is deliberately
// NOT the same as a system role, a position/title, a permission or authority.
const (
	StakeholderTypeOwner       = "owner"
	StakeholderTypeShareholder = "shareholder"
	StakeholderTypeBoardMember = "board_member"
	StakeholderTypeDirector    = "director"
	StakeholderTypeExecutive   = "executive"
	StakeholderTypeStaff       = "staff"
	StakeholderTypeConsultant  = "consultant"
	StakeholderTypePartner     = "partner"
	StakeholderTypeOther       = "other"
)

// StakeholderTypes is the initial catalogue of valid stakeholder types.
var StakeholderTypes = []string{
	StakeholderTypeOwner,
	StakeholderTypeShareholder,
	StakeholderTypeBoardMember,
	StakeholderTypeDirector,
	StakeholderTypeExecutive,
	StakeholderTypeStaff,
	StakeholderTypeConsultant,
	StakeholderTypePartner,
	StakeholderTypeOther,
}

// IsValidStakeholderType reports whether the value is in the initial catalogue.
func IsValidStakeholderType(value string) bool {
	for _, candidate := range StakeholderTypes {
		if candidate == value {
			return true
		}
	}
	return false
}

// Stakeholder is a generic relationship between one person (auth.User) and one
// organisation (eos.Entity), referenced by value.
//
// Foundation C1: this is the only table introduced by this release
// (ems_stakeholders). It does not create people, logins, roles or authority,
// and it has no foreign key into any existing or Education table.
type Stakeholder struct {
	gorm.Model
	Number             string     `json:"number" gorm:"size:36;unique;not null"`
	OrganisationNumber string     `json:"organisation_number" gorm:"size:36;not null;index"`
	UserNumber         string     `json:"user_number" gorm:"size:36;not null;index"`
	StakeholderType    string     `json:"stakeholder_type" gorm:"size:30;not null;index"`
	Title              string     `json:"title" gorm:"size:100"`
	SharePercentage    *float64   `json:"share_percentage" gorm:"type:decimal(6,3)"`
	Shares             *int64     `json:"shares"`
	ValidFrom          *time.Time `json:"valid_from"`
	ValidTo            *time.Time `json:"valid_to"`
	Status             string     `json:"status" gorm:"size:20;not null"`
	AddedByNumber      string     `json:"added_by_number" gorm:"size:36;not null"`
}

// TableName pins the table to ems_stakeholders.
func (Stakeholder) TableName() string {
	return "ems_stakeholders"
}
