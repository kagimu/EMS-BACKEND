package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateFeedingGroupUser(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.FeedingGroupUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	for _, user := range payload.UserNumber {

		checkExists, err := repositories.GetFeedingGroupUserByUserNumber(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Checking if user already added to a group error: " + err.Error()})
			return
		}

		if checkExists.Number != "" {
			continue
		}

		saveFeedingGroupUser := models.FeedingGroupUser{
			Number:             uuid.NewV4().String(),
			FeedingGroupNumber: payload.FeedingGroupNumber,
			Status:             payload.Status,
			UserNumber:         user,
			CampusNumber:       payload.CampusNumber,
		}

		if err := config.DB.Create(&saveFeedingGroupUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create feeding group user"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": "Feeding group user saved successfully"})
}

func FetchFeedingGroupUsers(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var feedingGroupUsers []models.FeedingGroupUser
	if err := config.DB.Preload("FeedingGroup").Preload("User").Find(&feedingGroupUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch feeding group users"})
		return
	}

	var response []responses.FeedingUserGroupResponse
	for _, feedingGroupUser := range feedingGroupUsers {
		response = append(response, responses.FeedingUserGroupResponse{
			Number:             feedingGroupUser.Number,
			UserNumber:         feedingGroupUser.UserNumber,
			Gender:             feedingGroupUser.User.Gender,
			Photo:              feedingGroupUser.User.Photo,
			Name:               feedingGroupUser.User.FirstName + " " + feedingGroupUser.User.SecondName + " " + feedingGroupUser.User.LastName,
			FeedingGroupNumber: feedingGroupUser.FeedingGroupNumber,
			FeedingGroupName:   feedingGroupUser.FeedingGroup.Name,
			Status:             feedingGroupUser.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
