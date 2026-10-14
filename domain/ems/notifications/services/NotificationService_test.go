package services

import (
	"errors"
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/models"
	"gorm.io/gorm"
)

func TestNotifyRejectsIncompleteInput(t *testing.T) {
	original := createNotification
	calls := 0
	createNotification = func(*models.Notification, ...*gorm.DB) error {
		calls++
		return nil
	}
	defer func() { createNotification = original }()

	for _, input := range []NotificationInput{
		{Type: "t", Title: "title"},
		{RecipientUserNumber: "u", Title: "title"},
		{RecipientUserNumber: "u", Type: "t"},
	} {
		if err := Notify(input); !errors.Is(err, ErrInvalidNotification) {
			t.Fatalf("Notify(%+v) = %v, want ErrInvalidNotification", input, err)
		}
	}
	if calls != 0 {
		t.Fatalf("invalid input reached the store %d times", calls)
	}
}

func TestNotifyCreatesUnreadInAppItem(t *testing.T) {
	original := createNotification
	var captured *models.Notification
	createNotification = func(notification *models.Notification, _ ...*gorm.DB) error {
		captured = notification
		return nil
	}
	defer func() { createNotification = original }()

	if err := Notify(NotificationInput{
		OrganisationNumber:  "org-1",
		RecipientUserNumber: "user-1",
		Type:                "approval.assigned",
		Title:               "Action required",
		Body:                "Please review",
		Payload:             map[string]any{"request": "r-1"},
	}); err != nil {
		t.Fatalf("Notify = %v", err)
	}
	if captured == nil {
		t.Fatal("Notify did not persist an item")
	}
	if captured.RecipientUserNumber != "user-1" || captured.Type != "approval.assigned" || captured.Title != "Action required" {
		t.Fatalf("unexpected item: %+v", captured)
	}
	if captured.ReadAt != nil {
		t.Fatal("a new notification must be unread")
	}
	if captured.Payload != `{"request":"r-1"}` {
		t.Fatalf("payload = %q", captured.Payload)
	}
}

func TestNotifyPropagatesStoreError(t *testing.T) {
	original := createNotification
	want := errors.New("store down")
	createNotification = func(*models.Notification, ...*gorm.DB) error {
		return want
	}
	defer func() { createNotification = original }()

	if err := Notify(NotificationInput{
		RecipientUserNumber: "user-1",
		Type:                "approval.assigned",
		Title:               "Action required",
	}); !errors.Is(err, want) {
		t.Fatalf("Notify = %v, want %v", err, want)
	}
}
