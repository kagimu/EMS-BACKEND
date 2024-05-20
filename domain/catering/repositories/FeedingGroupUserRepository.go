package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetFeedingGroupUserByGroupNumber(groupNumber string) ([]models.FeedingGroupUser, error) {

	db := config.DB

	var feedingGroupUsers []models.FeedingGroupUser

	if err := db.Preload("User").
		Where("feeding_group_number = ?", groupNumber).
		Find(&feedingGroupUsers).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingGroupUsers, nil
		}

		return feedingGroupUsers, err
	}

	return feedingGroupUsers, nil
}
