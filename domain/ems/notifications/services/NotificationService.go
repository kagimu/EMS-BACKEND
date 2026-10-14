package services

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/models"
	notificationRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/repositories"
	eosRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/repositories"
)

var (
	ErrNoActiveOrganisation = errors.New("no active organisation is set for this user")
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvalidNotification  = errors.New("notification recipient, type and title are required")
)

// createNotification is the in-app provider's persistence call used by Notify.
// It is a variable so a unit test can simulate a store failure.
var createNotification = notificationRepositories.CreateNotification

// NotificationInput is the input for creating one in-app notification item.
type NotificationInput struct {
	OrganisationNumber  string
	RecipientUserNumber string
	Type                string
	Title               string
	Body                string
	Payload             interface{}
}

// Notify writes one in-app notification item (Foundation G §11.2/§11.3).
//
// The recipient is determined by the caller from the ratified G-Q1 map; this
// function performs no recipient inference and no delivery beyond persisting
// the item. It never contacts an external provider. G3 does not yet call this
// from any foundation; that integration is G4.
func Notify(input NotificationInput) error {
	recipient := strings.TrimSpace(input.RecipientUserNumber)
	notificationType := strings.TrimSpace(input.Type)
	title := strings.TrimSpace(input.Title)
	if recipient == "" || notificationType == "" || title == "" {
		return ErrInvalidNotification
	}

	notification := models.Notification{
		OrganisationNumber:  strings.TrimSpace(input.OrganisationNumber),
		RecipientUserNumber: recipient,
		Type:                notificationType,
		Title:               title,
		Body:                input.Body,
		Payload:             toJSON(input.Payload),
	}
	return createNotification(&notification)
}

// ListRecipientNotifications lists the authenticated recipient's notifications
// within their active organisation, newest first. Both the organisation and the
// recipient are forced from the authenticated session, so a caller can never
// widen scope to another recipient or organisation (Foundation G §16/§20).
func ListRecipientNotifications(userNumber, notificationType string, unreadOnly bool, limit, offset int) ([]responses.NotificationResponse, int64, error) {
	organisationNumber, err := currentOrganisationNumber(userNumber)
	if err != nil {
		return nil, 0, err
	}

	filter := models.NotificationFilter{
		OrganisationNumber:  organisationNumber,
		RecipientUserNumber: userNumber,
		Type:                strings.TrimSpace(notificationType),
		UnreadOnly:          unreadOnly,
		Limit:               limit,
		Offset:              offset,
	}

	notifications, total, err := notificationRepositories.FindNotifications(filter)
	if err != nil {
		return nil, 0, err
	}

	result := make([]responses.NotificationResponse, 0, len(notifications))
	for _, notification := range notifications {
		result = append(result, toResponse(notification))
	}
	return result, total, nil
}

// GetRecipientNotification returns one of the authenticated recipient's
// notifications. A row addressed to another user or organisation is reported as
// not found so it cannot be probed (Foundation G §16/§20).
func GetRecipientNotification(userNumber string, id uint) (*responses.NotificationResponse, error) {
	organisationNumber, err := currentOrganisationNumber(userNumber)
	if err != nil {
		return nil, err
	}
	notification, err := notificationRepositories.FindNotificationForRecipient(organisationNumber, userNumber, id)
	if err != nil {
		return nil, err
	}
	if notification == nil {
		return nil, ErrNotificationNotFound
	}
	response := toResponse(*notification)
	return &response, nil
}

// MarkRecipientNotificationRead marks one of the authenticated recipient's
// notifications read. It changes only read_at and is idempotent: an already-read
// item is returned unchanged (Foundation G §23).
func MarkRecipientNotificationRead(userNumber string, id uint) (*responses.NotificationResponse, error) {
	organisationNumber, err := currentOrganisationNumber(userNumber)
	if err != nil {
		return nil, err
	}
	notification, err := notificationRepositories.FindNotificationForRecipient(organisationNumber, userNumber, id)
	if err != nil {
		return nil, err
	}
	if notification == nil {
		return nil, ErrNotificationNotFound
	}

	if notification.ReadAt == nil {
		now := time.Now()
		if err := notificationRepositories.MarkNotificationRead(notification.ID, now); err != nil {
			return nil, err
		}
		// Re-read so the response reflects the value actually persisted (the
		// store may hold a lower time precision than time.Now), keeping
		// repeated mark-read calls consistent and idempotent.
		persisted, err := notificationRepositories.FindNotificationForRecipient(organisationNumber, userNumber, id)
		if err != nil {
			return nil, err
		}
		if persisted != nil {
			notification = persisted
		} else {
			notification.ReadAt = &now
		}
	}

	response := toResponse(*notification)
	return &response, nil
}

func currentOrganisationNumber(userNumber string) (string, error) {
	active, err := eosRepositories.FindActiveUserEntityByUserNumber(userNumber)
	if err != nil {
		return "", err
	}
	if active.UserNumber == "" || active.EntityNumber == "" {
		return "", ErrNoActiveOrganisation
	}
	return active.EntityNumber, nil
}

func toJSON(v interface{}) string {
	if v == nil {
		return ""
	}
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

func toResponse(notification models.Notification) responses.NotificationResponse {
	response := responses.NotificationResponse{
		Number:              notification.ID,
		OrganisationNumber:  notification.OrganisationNumber,
		RecipientUserNumber: notification.RecipientUserNumber,
		Type:                notification.Type,
		Title:               notification.Title,
		Body:                notification.Body,
		ReadAt:              notification.ReadAt,
		CreatedAt:           notification.CreatedAt,
	}
	if notification.Payload != "" {
		response.Payload = json.RawMessage(notification.Payload)
	}
	return response
}
