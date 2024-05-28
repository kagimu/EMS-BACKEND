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

func GetFeedingSessionItemByFeedingSessionNumberAndMenuItemNumber(feedingSessionNumber string, menuItemNumber string) (models.FeedingSessionItem, error) {

	db := config.DB

	var feedingItem models.FeedingSessionItem

	if err := db.
		Where("menu_item_number = ?", menuItemNumber).
		Where("feeding_session_number = ?", feedingSessionNumber).
		First(&feedingItem).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingItem, nil
		}

		return feedingItem, err
	}

	return feedingItem, nil
}
