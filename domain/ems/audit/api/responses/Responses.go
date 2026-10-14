package responses

import (
	"encoding/json"
	"time"
)

// AuditLogResponse is the read model for one audit fact. It never carries
// licence keys, passwords, tokens or other secrets (Foundation G §20).
type AuditLogResponse struct {
	Number                 uint            `json:"number"`
	OrganisationNumber     string          `json:"organisation_number"`
	ActorUserNumber        string          `json:"actor_user_number"`
	Action                 string          `json:"action"`
	EntityType             string          `json:"entity_type"`
	EntityID               string          `json:"entity_id"`
	OldValues              json.RawMessage `json:"old_values,omitempty"`
	NewValues              json.RawMessage `json:"new_values,omitempty"`
	IP                     string          `json:"ip,omitempty"`
	UserAgent              string          `json:"user_agent,omitempty"`
	CorrelationID          string          `json:"correlation_id,omitempty"`
	WorkflowInstanceNumber string          `json:"workflow_instance_number,omitempty"`
	CreatedAt              time.Time       `json:"created_at"`
}
