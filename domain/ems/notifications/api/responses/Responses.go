package responses

import (
	"encoding/json"
	"time"
)

// NotificationResponse is the read model for one in-app notification. It is
// addressed to a single recipient and never carries secrets (Foundation G §20).
type NotificationResponse struct {
	Number              uint            `json:"number"`
	OrganisationNumber  string          `json:"organisation_number"`
	RecipientUserNumber string          `json:"recipient_user_number"`
	Type                string          `json:"type"`
	Title               string          `json:"title"`
	Body                string          `json:"body,omitempty"`
	Payload             json.RawMessage `json:"payload,omitempty"`
	ReadAt              *time.Time      `json:"read_at,omitempty"`
	CreatedAt           time.Time       `json:"created_at"`
}
