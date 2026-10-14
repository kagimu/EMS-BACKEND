package models

import "gorm.io/gorm"

// WorkflowDeadlineConfig is the Foundation E4 (ratified E4-D4) organisation-scoped
// configuration for the advisory approval-request deadline default.
//
// It is additive and minimal: one row per organisation holding the default
// deadline in HOURS (UTC-independent duration semantics, E4-D5). When no row
// exists the built-in default (7 days) applies. It never triggers any background
// behaviour (E4-D1) and never changes an existing request's persisted deadline.
type WorkflowDeadlineConfig struct {
	gorm.Model
	OrganisationNumber   string `json:"organisation_number" gorm:"size:36;uniqueIndex;not null"`
	DefaultDeadlineHours int    `json:"default_deadline_hours" gorm:"not null"`
}

// TableName pins the table to ems_workflow_deadline_configs.
func (WorkflowDeadlineConfig) TableName() string {
	return "ems_workflow_deadline_configs"
}
