package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchoolResults/models"
	"gorm.io/gorm"
)

type PreviousSchoolResultsRepository interface {
	FetchSchoolResultsByApplicationNumber(applicationNumber string, txn ...*gorm.DB) ([]models.PreviousSchoolResults, error)
}

type previousSchoolResultsRepository struct {
}

func (previousSchoolResultsRepository *previousSchoolResultsRepository) FetchSchoolResultsByApplicationNumber(applicationNumber string, txn ...*gorm.DB) ([]models.PreviousSchoolResults, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var previousSchoolResults []models.PreviousSchoolResults

	if err := db.Where("application_number = ?", applicationNumber).Find(&previousSchoolResults).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return previousSchoolResults, nil
		}

		return previousSchoolResults, err
	}

	return previousSchoolResults, nil
}

func NewPreviousSchoolResultsRepository() PreviousSchoolResultsRepository {
	return &previousSchoolResultsRepository{}
}
