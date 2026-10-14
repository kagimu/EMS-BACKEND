package repositories

import (
	"errors"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/models"
	"gorm.io/gorm"
)

// CreateNotification writes one in-app notification item. This is the in-app
// provider's persistence call (Foundation G §11.3). It accepts an optional
// transaction so a future G4 caller can keep the write with its originating
// action, without introducing a queue or background processor.
func CreateNotification(notification *models.Notification, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	return db.Create(notification).Error
}

// FindNotifications returns notifications matching the filter, newest first.
func FindNotifications(filter models.NotificationFilter) ([]models.Notification, int64, error) {
	db := config.DB.Model(&models.Notification{})
	if filter.OrganisationNumber != "" {
		db = db.Where("organisation_number = ?", filter.OrganisationNumber)
	}
	if filter.RecipientUserNumber != "" {
		db = db.Where("recipient_user_number = ?", filter.RecipientUserNumber)
	}
	if filter.Type != "" {
		db = db.Where("type = ?", filter.Type)
	}
	if filter.UnreadOnly {
		db = db.Where("read_at IS NULL")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var notifications []models.Notification
	db = db.Order("created_at desc")
	if filter.Limit > 0 {
		db = db.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		db = db.Offset(filter.Offset)
	}
	if err := db.Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// FindNotificationForRecipient returns one notification only when it belongs to
// the given organisation AND recipient. This is the recipient-isolation gate: a
// caller can never read another recipient's or another organisation's item
// (Foundation G §16/§20). Returns (nil, nil) when no such row exists.
func FindNotificationForRecipient(organisationNumber, recipientUserNumber string, id uint) (*models.Notification, error) {
	var notification models.Notification
	err := config.DB.
		Where("id = ? AND organisation_number = ? AND recipient_user_number = ?", id, organisationNumber, recipientUserNumber).
		First(&notification).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &notification, nil
}

// MarkNotificationRead sets only read_at on one notification. UpdateColumn is
// used deliberately: it touches no other column and does not alter the content
// or creation timestamp, so mark-read cannot mutate notification data
// (Foundation G §23). The caller must have already established that the row
// belongs to the acting recipient.
func MarkNotificationRead(id uint, readAt time.Time) error {
	return config.DB.Model(&models.Notification{}).
		Where("id = ?", id).
		UpdateColumn("read_at", readAt).Error
}
