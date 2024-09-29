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

	if err := db.Where("number = ?", schemeOfWorkSubTopicNumber).Find(&schemeOfWorkSubTopic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.SchemeOfWorkSubTopic{}, nil
		}

		return models.SchemeOfWorkSubTopic{}, err
	}

	return schemeOfWorkSubTopic, nil
}
