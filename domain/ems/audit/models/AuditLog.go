package models

import (
	"gorm.io/gorm"
)

// AuditLog is the Foundation G cross-cutting audit fact (ems_audit_logs).
//
// It is append-only: Foundation G provides creation and reading only, never
// update or deletion (Foundation G §6). The underlying business facts remain
// owned by Foundations A–F; this row references the authoritative entity by
// value (entity_type/entity_id) rather than replacing it (Foundation G §3.1,
// §5). Licensing history is NOT duplicated here: Foundation F's
// ems_licence_events remains the authoritative licensing fact source (§5.3,
// §17).
type AuditLog struct {
	gorm.Model
	OrganisationNumber string `json:"organisation_number" gorm:"size:36;not null;index:idx_ems_audit_organisation"`
	ActorUserNumber    string `json:"actor_user_number" gorm:"size:36;not null;index:idx_ems_audit_actor"`
	Action             string `json:"action" gorm:"size:64;not null;index:idx_ems_audit_action"`
	EntityType         string `json:"entity_type" gorm:"size:40;not null;index:idx_ems_audit_entity"`
	EntityID           string `json:"entity_id" gorm:"size:64;not null;index:idx_ems_audit_entity"`
	OldValues          string `json:"old_values" gorm:"type:text"`
	NewValues          string `json:"new_values" gorm:"type:text"`
	IP                 string `json:"ip" gorm:"size:45"`
	UserAgent          string `json:"user_agent" gorm:"size:255"`
	CorrelationID      string `json:"correlation_id" gorm:"size:64"`
	// WorkflowInstanceNumber is an optional reference to a Foundation E
	// workflow instance; it is never a foreign key.
	WorkflowInstanceNumber *string `json:"workflow_instance_number,omitempty" gorm:"size:36"`
	// CreatedAt is supplied by gorm.Model. There is deliberately no update or
	// delete path in the repository, service or API (§6).
}

// TableName pins the table to ems_audit_logs.
func (AuditLog) TableName() string {
	return "ems_audit_logs"
}

// AuditLogFilter narrows an audit query. Empty fields are ignored.
type AuditLogFilter struct {
	OrganisationNumber string
	ActorUserNumber    string
	Action             string
	EntityType         string
	EntityID           string
	Limit              int
	Offset             int
}
