package repository

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"gorm.io/gorm"
)

func SaveEnrollmentItem(enrollmentItem *models.EnrollmentItem, txn ...*gorm.DB) error {

	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}
	return db.Save(enrollmentItem).Error
}

func GetEnrollmentItemByID(id uint, txn ...*gorm.DB) (models.EnrollmentItem, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}
	var item models.EnrollmentItem
	if err := db.Preload("Item").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, nil
		}
		return item, err
	}
	return item, nil
}

func SaveRequirementClearanceLog(log models.RequirementClearanceLog, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}
	return db.Save(&log).Error
}

func FindRequirementClearanceLogs(enrollmentNumber, itemNumber, dateFrom, dateTo, staffNumber string) ([]models.RequirementClearanceLog, error) {
	db := config.DB
	query := db.Model(&models.RequirementClearanceLog{})

	if enrollmentNumber != "" {
		query = query.Where("enrollment_number = ?", enrollmentNumber)
	}
	if itemNumber != "" {
		query = query.Where("item_number = ?", itemNumber)
	}
	if dateFrom != "" && dateTo != "" {
		query = query.Where("date_cleared >= ? AND date_cleared <= ?", dateFrom, dateTo)
	} else if dateFrom != "" {
		query = query.Where("date_cleared >= ?", dateFrom)
	} else if dateTo != "" {
		query = query.Where("date_cleared <= ?", dateTo)
	}
	if staffNumber != "" {
		query = query.Where("cleared_by_staff_number = ?", staffNumber)
	}

	var logs []models.RequirementClearanceLog
	if err := query.Order("created_at desc").Find(&logs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return logs, nil
		}
		return logs, err
	}
	return logs, nil
}

// DeleteRequirementClearanceLogByID hard-deletes a clearance log record by its primary key.
func DeleteRequirementClearanceLogByID(id uint, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}
	return db.Unscoped().Delete(&models.RequirementClearanceLog{}, id).Error
}

// GetClearanceLogByID fetches a single clearance log record by ID (including soft-deleted check).
func GetClearanceLogByID(id uint) (models.RequirementClearanceLog, error) {
	var log models.RequirementClearanceLog
	if err := config.DB.First(&log, id).Error; err != nil {
		return log, err
	}
	return log, nil
}

// SumClearanceLogQuantityByItemID returns the total quantity_cleared across all remaining
// (non-deleted) logs for a given EnrollmentItemID. Used to recalculate the item's quantity
// after a log is deleted.
func SumClearanceLogQuantityByItemID(enrollmentItemID uint, txn ...*gorm.DB) (float64, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}
	var total float64
	err := db.Model(&models.RequirementClearanceLog{}).
		Where("enrollment_item_id = ?", enrollmentItemID).
		Select("COALESCE(SUM(quantity_cleared), 0)").
		Scan(&total).Error
	return total, err
}
