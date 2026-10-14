package decision

import (
	"log"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	notificationservices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/services"
)

// Event is one authoritative business event presented to the G4 notification
// decision layer. Only the fields relevant to the event's rule are read. The
// identifier fields are references by value; no business record is mutated.
type Event struct {
	Type               string
	OrganisationNumber string
	EntityType         string
	EntityID           string
	// Approval context.
	RequestNumber     string
	RequestedByNumber string
	Level             int
	// Foundation C context: the affected user.
	UserNumber string
	// Foundation F context: the licence/installation number (never the key).
	LicenceNumber string
	// SourceEventID is an optional deterministic key for an authoritative
	// source event (e.g. an ems_licence_events row id) used for best-effort
	// duplicate suppression. Empty where the source provides none.
	SourceEventID string
}

// notificationPayload is the structured data stored on the notification. It is
// a struct (not a map) so its JSON field order is deterministic, which the
// best-effort duplicate check relies on. It never contains secrets.
type notificationPayload struct {
	ActionRequired     bool   `json:"action_required"`
	Event              string `json:"event"`
	EntityType         string `json:"entity_type"`
	EntityID           string `json:"entity_id"`
	OrganisationNumber string `json:"organisation_number"`
	LicenceNumber      string `json:"licence_number,omitempty"`
	SourceEventID      string `json:"source_event_id,omitempty"`
}

// Notify applies the fixed G-Q1 map to one authoritative event. It resolves the
// approved recipient(s), deduplicates them by user number (D-GQ1-18), and
// creates one best-effort in-app notification per recipient.
//
// It is deliberately best-effort and returns nothing: a notification failure
// must never roll back or alter the authoritative business action (§19). An
// event with no rule, or a rule whose recipient cannot be resolved, produces no
// notification (FAIL CLOSED — no fallback recipient is invented, D-GQ1-4).
func Notify(event Event) {
	rule, mapped := rules[event.Type]
	if !mapped || event.OrganisationNumber == "" {
		return
	}

	recipients := resolveRecipients(rule.kinds, event)
	if len(recipients) == 0 {
		log.Printf("notification: no recipient resolved for %q (organisation %s)", event.Type, event.OrganisationNumber)
		return
	}

	seen := make(map[string]bool, len(recipients))
	for _, recipient := range recipients {
		if recipient == "" || seen[recipient] {
			// D-GQ1-18: one notification per user, even when the user
			// qualifies through several recipient categories.
			continue
		}
		seen[recipient] = true
		if event.SourceEventID != "" && alreadyNotified(event.OrganisationNumber, recipient, event.SourceEventID) {
			continue
		}
		createNotification(event, rule, recipient)
	}
}

// createNotification persists a single in-app notification, swallowing any
// failure so the caller's business result is never affected.
func createNotification(event Event, rule recipientRule, recipient string) {
	payload := notificationPayload{
		ActionRequired:     rule.actionRequired,
		Event:              event.Type,
		EntityType:         event.EntityType,
		EntityID:           event.EntityID,
		OrganisationNumber: event.OrganisationNumber,
		LicenceNumber:      event.LicenceNumber,
		SourceEventID:      event.SourceEventID,
	}
	input := notificationservices.NotificationInput{
		OrganisationNumber:  event.OrganisationNumber,
		RecipientUserNumber: recipient,
		Type:                event.Type,
		Title:               rule.title,
		Body:                rule.body,
		Payload:             payload,
	}
	if err := notificationservices.Notify(input); err != nil {
		log.Printf("notification: failed to create %q for %s: %v", event.Type, recipient, err)
	}
}

// alreadyNotified is a best-effort duplicate guard for events that carry a
// deterministic source event id. It performs a read against the existing
// ems_notifications table; no schema change is introduced. It is not
// concurrency-proof, but the G4 integration points are single-shot per
// authoritative action, so this is a safety net rather than the primary guard.
func alreadyNotified(organisationNumber, recipient, sourceEventID string) bool {
	needle := `%"source_event_id":"` + sourceEventID + `"%`
	var count int64
	err := config.DB.Table("ems_notifications").
		Where("deleted_at IS NULL AND organisation_number = ? AND recipient_user_number = ?", organisationNumber, recipient).
		Where("payload LIKE ?", needle).
		Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

// NotifyApprovalCreated notifies the current-level approvers that a request
// awaits their decision (D-GQ1-1, ACTION REQUIRED).
func NotifyApprovalCreated(organisationNumber, requestNumber string, level int) {
	Notify(Event{
		Type:               EventApprovalAssigned,
		OrganisationNumber: organisationNumber,
		EntityType:         "ems_approval_request",
		EntityID:           requestNumber,
		RequestNumber:      requestNumber,
		Level:              level,
	})
}

// NotifyApprovalAdvanced notifies the newly current-level approvers after an
// approval advances (D-GQ1-1, ACTION REQUIRED).
func NotifyApprovalAdvanced(organisationNumber, requestNumber string, level int) {
	Notify(Event{
		Type:               EventApprovalAdvanced,
		OrganisationNumber: organisationNumber,
		EntityType:         "ems_approval_request",
		EntityID:           requestNumber,
		RequestNumber:      requestNumber,
		Level:              level,
	})
}

// NotifyApprovalApproved informs the requester that the request was fully
// approved (D-GQ1-1/D-GQ1-10, INFORMATIONAL).
func NotifyApprovalApproved(organisationNumber, requestNumber, requesterUserNumber string) {
	Notify(Event{
		Type:               EventApprovalApproved,
		OrganisationNumber: organisationNumber,
		EntityType:         "ems_approval_request",
		EntityID:           requestNumber,
		RequestedByNumber:  requesterUserNumber,
	})
}

// NotifyApprovalRejected informs the requester that the request was rejected
// (D-GQ1-1/D-GQ1-8, INFORMATIONAL).
func NotifyApprovalRejected(organisationNumber, requestNumber, requesterUserNumber string) {
	Notify(Event{
		Type:               EventApprovalRejected,
		OrganisationNumber: organisationNumber,
		EntityType:         "ems_approval_request",
		EntityID:           requestNumber,
		RequestedByNumber:  requesterUserNumber,
	})
}

// NotifyApprovalReturned tells the requester the request was returned for
// correction (D-GQ1-1/D-GQ1-9, ACTION REQUIRED).
func NotifyApprovalReturned(organisationNumber, requestNumber, requesterUserNumber string) {
	Notify(Event{
		Type:               EventApprovalReturned,
		OrganisationNumber: organisationNumber,
		EntityType:         "ems_approval_request",
		EntityID:           requestNumber,
		RequestedByNumber:  requesterUserNumber,
	})
}

// NotifyStakeholderChanged notifies the affected stakeholder
// (D-GQ1-2/D-GQ1-11, INFORMATIONAL).
func NotifyStakeholderChanged(organisationNumber, entityID, affectedUserNumber string) {
	Notify(Event{
		Type:               EventStakeholderChanged,
		OrganisationNumber: organisationNumber,
		EntityType:         "ems_stakeholder",
		EntityID:           entityID,
		UserNumber:         affectedUserNumber,
	})
}

// NotifyAssignmentCreated notifies the assignee of a genuinely new
// responsibility (D-GQ1-13, ACTION REQUIRED).
func NotifyAssignmentCreated(organisationNumber, entityID, affectedUserNumber string) {
	Notify(Event{
		Type:               EventAssignmentAssigned,
		OrganisationNumber: organisationNumber,
		EntityType:         "ems_user_assignment",
		EntityID:           entityID,
		UserNumber:         affectedUserNumber,
	})
}

// NotifyAssignmentChanged notifies the assignee of a change to an existing
// assignment (D-GQ1-3/D-GQ1-12, INFORMATIONAL).
func NotifyAssignmentChanged(organisationNumber, entityID, affectedUserNumber string) {
	Notify(Event{
		Type:               EventAssignmentChanged,
		OrganisationNumber: organisationNumber,
		EntityType:         "ems_user_assignment",
		EntityID:           entityID,
		UserNumber:         affectedUserNumber,
	})
}

// NotifyLicenceFact applies the approved mapping to a Foundation F licensing
// fact. eventType is the exact LicenceEventType value and sourceEventID is the
// ems_licence_events row id (used for best-effort duplicate suppression). Only
// mapped licence events produce a notification; unmapped types are ignored.
func NotifyLicenceFact(eventType, organisationNumber, licenceNumber, sourceEventID string) {
	Notify(Event{
		Type:               eventType,
		OrganisationNumber: organisationNumber,
		EntityType:         "ems_licence",
		EntityID:           licenceNumber,
		LicenceNumber:      licenceNumber,
		SourceEventID:      sourceEventID,
	})
}
