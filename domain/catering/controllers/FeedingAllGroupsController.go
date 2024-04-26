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

func CreateFeedingAllGroups(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.FeedingPerGroupRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	feedingSessionNumber := payload.FeedingSessionNumber

	// checks in feedingSessionsPerDay if that meal is supposed to be served on that day of the week.
	var feedingSessionPerDay models.FeedingSessionPerDay
	if err := config.DB.Where("feeding_session_number = ? AND day = ?", feedingSessionNumber, payload.Day).Find(&feedingSessionPerDay).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding session for the specified day not found"})
		return
	}

	//gets all groups from feedingGroupPerDay by day of the week.
	var feedingGroupsPerDay models.FeedingGroupPerDay
	if err := config.DB.Where("day = ?", payload.Day).Find(&feedingGroupsPerDay).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding groups for the specified day not found"})
		return
	}

	// checks in feedingGroupSessions for each group fetched if that feeding group is entitled to that particular feeding session.
	var feedingGroupSession models.FeedingGroupSession
	if err := config.DB.Where("feeding_session_number = ?", feedingSessionNumber).Find(&feedingGroupSession).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding group session for the specified day not found"})
		return
	}

	saveFeedingForAll := models.FeedingAllGroups{
		Number:               uuid.NewV4().String(),
		FeedingSessionNumber: feedingSessionNumber,
		Day:                  payload.Day,
		UserNumber:           user.Number,
	}

	if err := config.DB.Create(&saveFeedingForAll).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save feeding for all groups"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feeding for all groups saved successfully"})
}

func FetchFeedingForAllGroups(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var feedingForAllGroups []models.FeedingAllGroups
	if err := config.DB.Find(&feedingForAllGroups).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch feeding for all groups"})
		return
	}

	var response []responses.FeedingsPerGroupResponse
	for _, feedingPerGroup := range feedingForAllGroups {
		response = append(response, responses.FeedingsPerGroupResponse{
			Number:               feedingPerGroup.Number,
			Day:                  feedingPerGroup.Day,
			FeedingSessionNumber: feedingPerGroup.FeedingSessionNumber,
			FeedingSessionName:   feedingPerGroup.FeedingSession.Name,
			FeedingSessionStatus: feedingPerGroup.FeedingSession.Status,
			FeedingSessionCode:   feedingPerGroup.FeedingSession.Code,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
