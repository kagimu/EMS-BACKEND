package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func GetFeedingGroupUserByGroupNumber(schoolNumber string, campusNumber string, groupNumber string) ([]models.FeedingGroupUser, error) {

	db := config.DB

	var feedingGroupUsers []models.FeedingGroupUser

	if err := db.Preload("User").Preload("Campus").Preload("FeedingGroup").
		Joins("JOIN feeding_groups ON feeding_group_users.feeding_group_number = feeding_groups.number").
		Scopes(filters.ByAnyField("feeding_groups.school_number", schoolNumber), filters.ByCampusNumber(campusNumber), filters.ByAnyField("feeding_group_number", groupNumber)).
		Find(&feedingGroupUsers).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingGroupUsers, nil
		}

		return feedingGroupUsers, err
	}

	return feedingGroupUsers, nil
}

func GetFeedingGroupUserByUserNumber(userNumber string) (models.FeedingGroupUser, error) {

	db := config.DB

	var feedingGroupUser models.FeedingGroupUser

	if err := db.Where("user_number", userNumber).First(&feedingGroupUser).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingGroupUser, nil
		}

		return feedingGroupUser, err
	}

	return feedingGroupUser, nil
}

func AvailableUsers(schoolNumber string, status string, userType string) ([]userModel.User, error) {

	db := config.DB

	var users = make([]userModel.User, 0)

	if err := db.Scopes(filters.ByAnyField("users.status", status), filters.ByType(userType), notInGroupUsers()).Find(&users).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return users, nil
		}

		return users, err
	}

	return users, nil
}

func notInGroupUsers() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins("LEFT JOIN feeding_group_users ON users.number = feeding_group_users.user_number").
			Where("feeding_group_users.user_number IS NULL")
	}
}
