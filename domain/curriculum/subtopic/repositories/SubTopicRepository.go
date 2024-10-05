package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	subTopicModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	"gorm.io/gorm"
)

func GetSubTopicByNumber(subTopicNumber string, tx ...*gorm.DB) (subTopicModel.SubTopic, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var subTopic subTopicModel.SubTopic
	if err := db.Where("number = ?", subTopicNumber).First(&subTopic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return subTopicModel.SubTopic{}, nil
		}

		return subTopicModel.SubTopic{}, err
	}

	return subTopic, nil
}
