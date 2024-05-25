package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func GetFeedingGroupByNumber(feedingGroupNumber string) (models.FeedingGroup, error) {

	var feedingGroup models.FeedingGroup

	db := config.DB

	if err := db.Where("number = ?", feedingGroupNumber).First(&feedingGroup).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingGroup, nil
		}

		return feedingGroup, err
	}

	return feedingGroup, nil

}

func GetAllFeedingGroups(schoolNumber string, status string) ([]models.FeedingGroup, error) {

	var feedingGroups []models.FeedingGroup

	db := config.DB

	if err := db.Scopes(
		filters.BySchoolNumber(schoolNumber),
		filters.ByStatus(status)).
		Find(&feedingGroups).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingGroups, nil
		}

		return feedingGroups, err
	}

	return feedingGroups, nil
}

func GetFeedingGroupByName(name string, schoolNumber string) (models.FeedingGroup, error) {

	var feedingGroup models.FeedingGroup

	db := config.DB

	if err := db.
		Where("name = ?", name).
		Where("school_number = ?", schoolNumber).
		First(&feedingGroup).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingGroup, nil
		}

		return feedingGroup, err
	}

	return feedingGroup, nil
}

func SaveFeedingGroup(feedingGroup models.FeedingGroup) error {
	
	if err := config.DB.Save(&feedingGroup).Error; err != nil {
		return err
	}

	return nil
}
