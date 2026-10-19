package models

import (
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"gorm.io/gorm"
)

const (
	ReportingStatusActive   = constants.StatusActive
	ReportingStatusInactive = constants.StatusInactive
	ReportingStatusEnded    = "ENDED"
)

// ReportingLine records that one user assignment reports to another user
// assignment within the same organisation.
//
// It is intentionally attached to assignments, not stakeholders or the
// organisation-unit hierarchy: an assignment is where a person operates and in
// what capacity, while a stakeholder row is a governance/ownership relationship
// and the unit parent hierarchy is structural only. Both ends are stored by
// value (no foreign keys), preserving history even if an assignment is later
// ended; cross-organisation and cycle rules are enforced by the service.
type ReportingLine struct {
	gorm.Model
	Number                      string     `json:"number" gorm:"size:36;unique;not null"`
	OrganisationNumber          string     `json:"organisation_number" gorm:"size:36;not null;index;index:idx_ems_reporting_lines_org_status,priority:1"`
	SubordinateAssignmentNumber string     `json:"subordinate_assignment_number" gorm:"size:36;not null;index"`
	ManagerAssignmentNumber     string     `json:"manager_assignment_number" gorm:"size:36;not null;index"`
	ValidFrom                   *time.Time `json:"valid_from"`
	ValidTo                     *time.Time `json:"valid_to"`
	Status                      string     `json:"status" gorm:"size:20;not null;index:idx_ems_reporting_lines_org_status,priority:2"`
	AddedByNumber               string     `json:"added_by_number" gorm:"size:36;not null"`
}

// TableName pins the table to ems_reporting_lines.
func (ReportingLine) TableName() string {
	return "ems_reporting_lines"
}
