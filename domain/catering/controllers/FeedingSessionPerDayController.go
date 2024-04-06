package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateFeedingSessionPerDay(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.FeedingSessionPerDayRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	saveFeedingSessionPerDay := models.FeedingGroupPerDay{
		Number:             uuid.NewV4().String(),
		FeedingGroupNumber: payload.FeedingSessionNumber,
		Day:                payload.Day,
		UserNumber:         user.Number,
	}

	if err := config.DB.Create(&saveFeedingSessionPerDay).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save feeding session on a specified day"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": saveFeedingSessionPerDay})
}

func FetchFeedingSessionsPerDay(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var feedingSessionsPerDay []models.FeedingSessionPerDay
	if err := config.DB.Find(&feedingSessionsPerDay).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch feeding sessions per day"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": feedingSessionsPerDay})
}
