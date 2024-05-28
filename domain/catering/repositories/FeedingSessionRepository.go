package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetFeedingSessionByNumber(feedingSession string) (models.FeedingSession, error) {

	var session models.FeedingSession

	db := config.DB

	if err := db.Where("number = ?", feedingSession).First(&session).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return session, nil
		}

		return session, err
	}

	return session, nil
}

func GetFeedingSessionByNameAndSchool(name string, schoolNumber string) (models.FeedingSession, error) {
	var session models.FeedingSession

	db := config.DB

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

func SaveFeedingSession(feedingSession models.FeedingSession) error {

	if err := config.DB.Save(&feedingSession).Error; err != nil {
		return err
	}

	return nil
}

func GetFeedingSessionByCodeAndSchool(code uint64, schoolNumber string) (models.FeedingSession, error) {
	var session models.FeedingSession

	db := config.DB

	if err := db.
		Where("school_number = ?", schoolNumber).
		Where("code = ?", code).
		First(&session).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return session, nil
		}

		return session, err
	}

	return session, nil
}

func GetFeedingSessionsBySchoolNumber(schoolNumber string) ([]models.FeedingSession, error) {
	var sessions []models.FeedingSession

	db := config.DB

	if err := db.Where("school_number = ?", schoolNumber).Find(&sessions).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sessions, nil
		}

		return sessions, err
	}

	return sessions, nil
}
