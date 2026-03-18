package repositories

import (
	"errors"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	invoiceModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
)

func SavePrepaidDeposit(deposit invoiceModels.PrepaidDeposit, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	return db.Create(&deposit).Error
}

func FindPrepaidBalanceByUserNumber(userNumber string, tx ...*gorm.DB) (invoiceModels.PrepaidBalance, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	var b invoiceModels.PrepaidBalance
	err := db.Where("user_number = ?", userNumber).First(&b).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invoiceModels.PrepaidBalance{UserNumber: userNumber, Balance: 0, UpdatedAt: time.Now()}, nil
		}
		return b, err
	}
	return b, nil
}

// AddToPrepaidBalance adds amount to the user's prepaid balance (creates row if not exists).
func AddToPrepaidBalance(userNumber string, amount float64, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	return db.Model(&invoiceModels.PrepaidBalance{}).
		Where("user_number = ?", userNumber).
		Updates(map[string]interface{}{
			"balance":    gorm.Expr("COALESCE(balance, 0) + ?", amount),
			"updated_at": time.Now(),
		}).Error
}

// EnsurePrepaidBalanceRow creates a prepaid_balances row for user_number if it does not exist, then adds amount.
func EnsurePrepaidBalanceRow(userNumber string, amount float64, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	var exists int64
	if err := db.Model(&invoiceModels.PrepaidBalance{}).Where("user_number = ?", userNumber).Count(&exists).Error; err != nil {
		return err
	}
	if exists == 0 {
		if err := db.Create(&invoiceModels.PrepaidBalance{
			UserNumber: userNumber,
			Balance:    amount,
			UpdatedAt:  time.Now(),
		}).Error; err != nil {
			return err
		}
		return nil
	}
	return AddToPrepaidBalance(userNumber, amount, tx...)
}

func FindDepositsByUserNumber(userNumber string, tx ...*gorm.DB) ([]invoiceModels.PrepaidDeposit, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	var list []invoiceModels.PrepaidDeposit
	err := db.Where("user_number = ?", userNumber).Order("created_at DESC").Find(&list).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return list, nil
}
