package models

import (
	"time"

	"gorm.io/gorm"
)

// Assignment types. Open set: an assignment records WHERE a person operates and
// in what capacity. It is deliberately NOT a system role, a stakeholder type, a
// permission, or authority.
const (
	AssignmentTypeMember  = "member"
	AssignmentTypeManager = "manager"
	AssignmentTypeHead    = "head"
	AssignmentTypeStaff   = "staff"
)

// AssignmentTypes is the initial catalogue of valid assignment types.
var AssignmentTypes = []string{
	AssignmentTypeMember,
	AssignmentTypeManager,
	AssignmentTypeHead,
	AssignmentTypeStaff,
}

// IsValidAssignmentType reports whether the value is in the initial catalogue.
func IsValidAssignmentType(value string) bool {
	for _, candidate := range AssignmentTypes {
		if candidate == value {
			return true
		}
	}
	return false
}

// UserAssignment records where a person operates within an organisation's unit
// hierarchy: organisation (eos.Entity) + user (auth.User) + organisation unit
// (ems_organisation_units), all referenced by value.
//
// Foundation C3: this is the only table introduced by this release
// (ems_user_assignments). It grants no permission or authority, has no foreign
// key into any existing or Education table, and is never read by Education.
type UserAssignment struct {
	gorm.Model
	Number                 string     `json:"number" gorm:"size:36;unique;not null"`
	OrganisationNumber     string     `json:"organisation_number" gorm:"size:36;not null;index;index:idx_ems_user_assignments_org_status,priority:1"`
	UserNumber             string     `json:"user_number" gorm:"size:36;not null;index"`
	OrganisationUnitNumber string     `json:"organisation_unit_number" gorm:"size:36;not null;index"`
	AssignmentType         string     `json:"assignment_type" gorm:"size:30;not null;index"`
	PositionTitle          string     `json:"position_title" gorm:"size:100"`
	ValidFrom              *time.Time `json:"valid_from"`
	ValidTo                *time.Time `json:"valid_to"`
	Status                 string     `json:"status" gorm:"size:20;not null;index:idx_ems_user_assignments_org_status,priority:2"`
	AddedByNumber          string     `json:"added_by_number" gorm:"size:36;not null"`
}

// TableName pins the table to ems_user_assignments.
func (UserAssignment) TableName() string {
	return "ems_user_assignments"
}
