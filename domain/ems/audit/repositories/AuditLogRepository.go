package repositories

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
	"gorm.io/gorm"
)

// CreateAuditLog appends an audit fact. This repository intentionally exposes
// no update or delete operation: audit records are immutable (Foundation G §6).
func CreateAuditLog(log *models.AuditLog, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	return db.Create(log).Error
}

// FindAuditLogs returns audit records matching the filter, newest first.
func FindAuditLogs(filter models.AuditLogFilter) ([]models.AuditLog, int64, error) {
	db := config.DB.Model(&models.AuditLog{})
	if filter.OrganisationNumber != "" {
		db = db.Where("organisation_number = ?", filter.OrganisationNumber)
	}
	if filter.ActorUserNumber != "" {
		db = db.Where("actor_user_number = ?", filter.ActorUserNumber)
	}
	if filter.Action != "" {
		db = db.Where("action = ?", filter.Action)
	}
	if filter.EntityType != "" {
		db = db.Where("entity_type = ?", filter.EntityType)
	}
	if filter.EntityID != "" {
		db = db.Where("entity_id = ?", filter.EntityID)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []models.AuditLog
	db = db.Order("created_at desc")
	if filter.Limit > 0 {
		db = db.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		db = db.Offset(filter.Offset)
	}
	if err := db.Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// FindAuditLogByID returns one record, or (nil, nil) when it does not exist.
func FindAuditLogByID(id uint) (*models.AuditLog, error) {
	var log models.AuditLog
	if err := config.DB.Where("id = ?", id).First(&log).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}
