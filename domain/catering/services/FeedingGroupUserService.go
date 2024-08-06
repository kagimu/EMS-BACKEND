package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreateFeedingGroupUser(userNumber, campusNumber, feedingGroupNumber, status string, txn *gorm.DB) error {

	checkExists, err := repositories.GetFeedingGroupUserByUserNumber(userNumber)
	if err != nil {
		return errors.New("Checking if user already added to a group error: " + err.Error())
	}

	if checkExists.Number != "" {
		return errors.New("user already added to a group")
	}

	saveFeedingGroupUser := models.FeedingGroupUser{
		Number:             uuid.NewV4().String(),
		FeedingGroupNumber: feedingGroupNumber,
		Status:             status,
		UserNumber:         userNumber,
		CampusNumber:       campusNumber,
	}

	if err := repositories.CreateFeedingGroupUser(saveFeedingGroupUser, txn).Error; err != nil {
		return errors.New("failed to create feeding group user ")
	}

	return nil

}
