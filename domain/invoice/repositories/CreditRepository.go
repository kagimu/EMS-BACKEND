package repositories

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	invoiceModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
)

func SaveCredit(credit invoiceModels.Credit, tx ...*gorm.DB) (invoiceModels.Credit, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Save(&credit).Error; err != nil {
		return credit, err
	}

	return credit, nil
}

func FindCreditByNumber(creditNumber string, tx ...*gorm.DB) (invoiceModels.Credit, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var credit invoiceModels.Credit
	if err := db.
		Preload("User").
		Preload("AddedBy").
		Where("number = ?", creditNumber).First(&credit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoiceModels.Credit{}, nil
		}
		return invoiceModels.Credit{}, err
	}

	return credit, nil
}

func FindCreditsByUserNumber(userNumber string, tx ...*gorm.DB) ([]invoiceModels.Credit, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var credits []invoiceModels.Credit
	if err := db.
		Where("user_number = ?", userNumber).
		Order("date desc").
		Find(&credits).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return credits, nil
		}
		return credits, err
	}

	return credits, nil
}

func GetCreditByClientReference(clientReference, campusNumber string, tx ...*gorm.DB) (invoiceModels.Credit, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var credit invoiceModels.Credit
	if err := db.
		Where("client_reference = ? AND campus_number = ?", clientReference, campusNumber).First(&credit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoiceModels.Credit{}, nil
		}
		return invoiceModels.Credit{}, err
	}

	return credit, nil
}
