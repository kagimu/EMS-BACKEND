package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input- " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tx := config.DB.Begin()
	for _, day := range payload.Day {

		var duplicateFeedingSession models.FeedingSessionPerDay
		if err := config.DB.Where("day = ? AND feeding_session_number = ?", day, payload.FeedingSessionNumber).First(&duplicateFeedingSession).Error; err == nil {
			continue
		}

		saveFeedingSessionPerDay := models.FeedingSessionPerDay{
			Number:               uuid.NewV4().String(),
			FeedingSessionNumber: payload.FeedingSessionNumber,
			Day:                  day,
			UserNumber:           user.Number,
		}

		if err := tx.Create(&saveFeedingSessionPerDay).Error; err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}

	if err = tx.Commit().Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to attach days to feeding session"})
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "Feeding session on the given day has been created"})
}

func FetchFeedingSessionsPerDay(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var feedingSessionsPerDay []models.FeedingSessionPerDay
	if err := config.DB.Preload("FeedingSession").Find(&feedingSessionsPerDay).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch feeding sessions per day"})
		return
	}

	var response []responses.FeedingSessionPerDayResponse
	for _, feedingSession := range feedingSessionsPerDay {
		response = append(response, responses.FeedingSessionPerDayResponse{
			Number:               feedingSession.Number,
			FeedingSessionNumber: feedingSession.FeedingSessionNumber,
			Day:                  feedingSession.Day,
			FeedingSessionName:   feedingSession.FeedingSession.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
