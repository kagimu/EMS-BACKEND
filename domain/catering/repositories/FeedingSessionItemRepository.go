package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetFeedingSessionItems(feedingSessionNumber string, tx ...*gorm.DB) ([]models.FeedingSessionItem, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

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

func GetFeedingSessionItemByFeedingSessionNumberAndMenuItemNumber(feedingSessionNumber, menuItemNumber string, tx ...*gorm.DB) (models.FeedingSessionItem, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

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

func SaveFeedingSessionItem(feedingSessionItem models.FeedingSessionItem, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Save(&feedingSessionItem).Error

}

func RemoveFeedingSessionItem(feedingSessionItem models.FeedingSessionItem, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Unscoped().Delete(&feedingSessionItem).Error

}
