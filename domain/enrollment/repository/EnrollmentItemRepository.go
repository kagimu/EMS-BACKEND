package repository

import (
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
