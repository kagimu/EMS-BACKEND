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

func CreateFeedingPerGroup(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.FeedingPerGroupRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dayError, isValid := commons.CheckDay(payload.Day)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"message": dayError})
		return
	}

	feedingSessionNumber := payload.FeedingSessionNumber
	var feedingSession models.FeedingSession
	if err := config.DB.Where("number = ?", feedingSessionNumber).Find(&feedingSession).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding session not found"})
		return
	}

	feedingGroupNumber := payload.FeedingGroupNumber
	var feedingGroup models.FeedingGroup
	if err := config.DB.Where("number = ?", feedingGroupNumber).Find(&feedingGroup).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding group not found"})
		return
	}

	saveFeedingPerGroup := models.FeedingPerGroup{
		Number:               uuid.NewV4().String(),
		FeedingGroupNumber:   payload.FeedingGroupNumber,
		FeedingSessionNumber: payload.FeedingSessionNumber,
		Day:                  payload.Day,
		UserNumber:           user.Number,
	}

	var feedingSessionPerDay models.FeedingSessionPerDay
	if err := config.DB.Where("feeding_session_number = ? AND day = ?", feedingSessionNumber, payload.Day).Find(&feedingSessionPerDay).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding session for the specified day not found"})
		return
	}

	var feedingGroupPerDay models.FeedingGroupPerDay
	if err := config.DB.Where("feeding_group_number = ? AND day = ?", feedingGroupNumber, payload.Day).Find(&feedingGroupPerDay).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding group for the specified day not found"})
		return
	}

	var feedingSessionPerGroup models.FeedingGroupSession
	if err := config.DB.Where("feeding_session_number = ? AND feeding_group_number = ?", feedingSessionNumber, feedingGroupNumber).Find(&feedingSessionPerGroup).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding session for the specified group not found"})
		return
	}

	if err := config.DB.Create(&saveFeedingPerGroup).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save feeding per group", "err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feeding per group saved successfully"})
}

func FetchFeedingsPerGroup(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var feedingsPerGroup []models.FeedingPerGroup
	if err := config.DB.Preload("FeedingGroup").Preload("FeedingSession").Find(&feedingsPerGroup).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch feedings per group"})
		return
	}

	var response []responses.FeedingsPerGroupResponse
	for _, feedingPerGroup := range feedingsPerGroup {
		response = append(response, responses.FeedingsPerGroupResponse{
			Number:               feedingPerGroup.Number,
			GroupName:            feedingPerGroup.FeedingGroup.Name,
			GroupStatus:          feedingPerGroup.FeedingGroup.Status,
			Day:                  feedingPerGroup.Day,
			FeedingSessionNumber: feedingPerGroup.FeedingSessionNumber,
			FeedingSessionName:   feedingPerGroup.FeedingSession.Name,
			FeedingSessionStatus: feedingPerGroup.FeedingSession.Status,
			FeedingSessionCode:   feedingPerGroup.FeedingSession.Code,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
