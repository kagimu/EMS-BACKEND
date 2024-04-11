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

func CreateFeedingSession(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.FeedingSessionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveFeedingSession := models.FeedingSession{
		Number:       uuid.NewV4().String(),
		Name:         payload.Name,
		Status:       payload.Status,
		CampusNumber: payload.CampusNumber,
		Code:         payload.Code,
		UserNumber:   user.Number,
	}

	var duplicateFeedingSession models.FeedingSession
	if err := config.DB.Where("name = ?", payload.Name).First(&duplicateFeedingSession).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding session name already exists"})
		return
	}

	if err := config.DB.Create(&saveFeedingSession).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save feeding session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Feeding session saved successfully"})
}

func FetchFeedingSessions(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var feedingSessions []models.FeedingSession
	if err := config.DB.Find(&feedingSessions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch feeding sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": feedingSessions})
}
