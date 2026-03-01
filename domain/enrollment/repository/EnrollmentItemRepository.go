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
	if err := db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, nil
		}
		return item, err
	}
	return item, nil
}
