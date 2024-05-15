package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetFeedingSessionItems(feedingSessionNumber string) ([]models.FeedingSessionItem, error) {

	db := config.DB

	var feedingSessionItems []models.FeedingSessionItem

	if err := db.
		Where("feeding_session_number = ? ", feedingSessionNumber).
		Preload("MenuItem").
		Find(&feedingSessionItems).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingSessionItems, nil
		}

		return feedingSessionItems, err
	}

	return feedingSessionItems, nil
}
