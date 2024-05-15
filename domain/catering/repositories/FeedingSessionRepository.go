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
