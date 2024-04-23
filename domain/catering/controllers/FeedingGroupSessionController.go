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

func CreateFeedingGroupSession(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.FeedingGroupSessionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveFeedingGroupSession := models.FeedingGroupSession{
		Number:               uuid.NewV4().String(),
		FeedingGroupNumber:   payload.FeedingGroupNumber,
		Status:               payload.Status,
		FeedingSessionNumber: payload.FeedingSessionNumber,
	}

	if err := config.DB.Create(&saveFeedingGroupSession).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create feeding group session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Feeding group session saved successfully"})
}

func FetchFeedingGroupSessions(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var feedingGroupSessions []models.FeedingGroupSession
	if err := config.DB.Preload("FeedingGroup").Preload("FeedingSession").Find(&feedingGroupSessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch feeding group users"})
		return
	}

	var response []responses.FeedingGroupSessionResponse
	for _, feedingGroupSession := range feedingGroupSessions {
		response = append(response, responses.FeedingGroupSessionResponse{
			Number:               feedingGroupSession.Number,
			FeedingGroupNumber:   feedingGroupSession.FeedingGroupNumber,
			Status:               feedingGroupSession.Status,
			FeedingSessionNumber: feedingGroupSession.FeedingSessionNumber,
			FeedingGroupName:     feedingGroupSession.FeedingGroup.Name,
			FeedingSessionName:   feedingGroupSession.FeedingSession.Name,
			FeedingSessionStatus: feedingGroupSession.FeedingSession.Status,
			FeedingSessionCampus: feedingGroupSession.FeedingGroup.Campus.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
