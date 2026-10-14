package models

import (
	"time"

	"gorm.io/gorm"
)

// Notification is a Foundation G in-app notification item (ems_notifications).
//
// Notifications are recipient-scoped: a row is addressed to exactly one user in
// one organisation (Foundation G §10.1/§16). Foundation G only ever creates and
// reads them, and a mark-read operation changes only ReadAt (§23). There is
// deliberately no external delivery, scheduler, worker, queue or background
// processor (§14/§21).
type Notification struct {
	gorm.Model
	OrganisationNumber  string     `json:"organisation_number" gorm:"size:36;not null;index:idx_ems_notification_organisation"`
	RecipientUserNumber string     `json:"recipient_user_number" gorm:"size:36;not null;index:idx_ems_notification_recipient"`
	Type                string     `json:"type" gorm:"size:64;not null;index:idx_ems_notification_type"`
	Title               string     `json:"title" gorm:"size:255;not null"`
	Body                string     `json:"body" gorm:"type:text"`
	Payload             string     `json:"payload" gorm:"type:text"`
	ReadAt              *time.Time `json:"read_at"`
}

// TableName pins the table to ems_notifications.
func (Notification) TableName() string {
	return "ems_notifications"
}

// NotificationFilter narrows a notification query. Empty fields are ignored.
type NotificationFilter struct {
	OrganisationNumber  string
	RecipientUserNumber string
	Type                string
	UnreadOnly          bool
	Limit               int
	Offset              int
}
