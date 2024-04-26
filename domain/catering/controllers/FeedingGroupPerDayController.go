package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateFeedingGroupPerDay(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.FeedingGroupPerDayRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}

	saveFeedingGroupPerDay := models.FeedingGroupPerDay{
		Number:             uuid.NewV4().String(),
		FeedingGroupNumber: payload.FeedingGroupNumber,
		Day:                payload.Day,
		UserNumber:         user.Number,
	}

	var duplicateFeedingGroup models.FeedingGroupPerDay
	if err := config.DB.Where("day = ? AND feeding_group_number = ?", payload.Day, payload.FeedingGroupNumber).First(&duplicateFeedingGroup).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding group on the given day already exists"})
		return
	}

	if err := config.DB.Create(&saveFeedingGroupPerDay).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save feeding group on a specified day"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feeding group per given day created successfully"})
}

func FetchFeedingGroupsPerDay(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var feedingGroupsPerDay []models.FeedingGroupPerDay
	if err := config.DB.Preload("FeedingGroup").Find(&feedingGroupsPerDay).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch feeding groups per day"})
		return
	}

	var response []responses.FeedingGroupPerDayResponse
	for _, feedingGroupPerDay := range feedingGroupsPerDay {
		response = append(response, responses.FeedingGroupPerDayResponse{
			Number:             feedingGroupPerDay.Number,
			FeedingGroupNumber: feedingGroupPerDay.FeedingGroupNumber,
			FeedingGroupName:   feedingGroupPerDay.FeedingGroup.Name,
			Day:                feedingGroupPerDay.Day,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
