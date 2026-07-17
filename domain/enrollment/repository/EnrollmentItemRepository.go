package repository

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"gorm.io/gorm"
)

func SaveEnrollmentItem(enrollmentItem models.EnrollmentItem, txn ...*gorm.DB) error {

	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}
	return db.Save(&enrollmentItem).Error
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

