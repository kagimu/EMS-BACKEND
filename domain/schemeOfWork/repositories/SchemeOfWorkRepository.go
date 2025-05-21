package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/models"
	"gorm.io/gorm"
)

func GetSchemeOfWorkSubTopicByNumber(schemeOfWorkSubTopicNumber string, tx ...*gorm.DB) (models.SchemeOfWorkSubTopic, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var schemeOfWorkSubTopic models.SchemeOfWorkSubTopic

	if err := db.Where("number = ?", schemeOfWorkSubTopicNumber).First(&schemeOfWorkSubTopic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.SchemeOfWorkSubTopic{}, nil
		}

		return models.SchemeOfWorkSubTopic{}, err
	}

	return schemeOfWorkSubTopic, nil
}

func SaveSchemeOfWork(schemeOfWork models.SchemeOfWork, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Save(&schemeOfWork).Error; err != nil {
		return err
	}

	return nil
}

func GetSchemeOfWorkByNumber(schemeOfWorkNumber string, tx ...*gorm.DB) (models.SchemeOfWork, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var schemeOfWorkSubTopic models.SchemeOfWork

	if err := db.Where("number = ?", schemeOfWorkNumber).First(&schemeOfWorkSubTopic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.SchemeOfWork{}, nil
		}

		return models.SchemeOfWork{}, err
	}

	return schemeOfWorkSubTopic, nil
}
