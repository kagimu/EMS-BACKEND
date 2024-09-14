package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"gorm.io/gorm"
)

func GetStatementByNumber(statementNumber string, tx ...*gorm.DB) (models.Statement, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	var statement models.Statement

	if err := db.Where("number = ?", statementNumber).First(&statement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Statement{}, nil
		}

		return models.Statement{}, err
	}

	return statement, nil
}

func SaveStatement(statement models.Statement, tx ...*gorm.DB) error {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Save(&statement).Error
}

func GetRecentStatementByUserNumber(userNumber string, tx ...*gorm.DB) (models.Statement, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	var statement models.Statement

	if err := db.Where("user_number = ?", userNumber).First(&statement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Statement{}, nil
		}

		return models.Statement{}, err
	}

	return statement, nil
}
