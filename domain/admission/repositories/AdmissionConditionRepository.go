package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"gorm.io/gorm"
)

func FindAllAdmissionItem(admissionNumber string, txn ...*gorm.DB) ([]models.AdmissionCondition, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var admissionConditions []models.AdmissionCondition

	if err := db.
		Joins("JOIN items On items.number = admission_conditions.item_number").
		Where("admission_number = ? AND items.mandatory = ?", admissionNumber, true).Find(&admissionConditions).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return admissionConditions, nil
		}

		return admissionConditions, err
	}

	return admissionConditions, nil
}
