package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	uuid "github.com/satori/go.uuid"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
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

	schoolNumber := c.Query("school_number")
	groupNumber := c.Query("group_number")
	campusNumber := c.Query("campus_number")

	response, err := services.MealUsers(schoolNumber, campusNumber, groupNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func FetchAvailableUsers(context *gin.Context) {
	_, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	schoolNumber := context.Query("school_number")
	userType := context.Query("type")

	response, err := services.AvailableUser(schoolNumber, userType)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": response})

}
