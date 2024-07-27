package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func GetFeedingSessionByNumberAndSchoolNumber(feedingSession, schoolNumber string, tx ...*gorm.DB) (models.FeedingSession, error) {

	var session models.FeedingSession

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Where("number = ?", feedingSession).Scopes(filters.BySchoolNumber(schoolNumber)).First(&session).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return session, nil
		}

		return session, err
	}

	return session, nil
}

func GetFeedingSessionByNameAndSchool(name, schoolNumber string, tx ...*gorm.DB) (models.FeedingSession, error) {
	var session models.FeedingSession

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.
		Where("school_number = ?", schoolNumber).
		Where("name = ?", name).
		First(&session).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return session, nil
		}

		return session, err
	}

	return session, nil
}

func SaveFeedingSession(feedingSession models.FeedingSession, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Save(&feedingSession).Error; err != nil {
		return err
	}

	return nil
}

func GetFeedingSessionsBySchoolNumber(schoolNumber string, tx ...*gorm.DB) ([]models.FeedingSession, error) {
	var sessions []models.FeedingSession

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Where("school_number = ?", schoolNumber).Find(&sessions).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sessions, nil
		}

		return sessions, err
	}

	return sessions, nil
}

func DeleteFeedingSession(feedingSession models.FeedingSession, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Delete(&feedingSession).Error; err != nil {
		return err
	}

	return nil

}
