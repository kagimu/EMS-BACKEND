package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateFeedingGroup(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.FeedingGroupRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	campusNumber := payload.CampusNumber
	var campus campusModels.Campus
	if err := config.DB.First(&campus, "number = ?", campusNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No campus found with the specified number"})
		return
	}

	saveFeedingGroup := models.FeedingGroup{
		Number:       uuid.NewV4().String(),
		Name:         payload.Name,
		CampusNumber: campusNumber,
		Status:       payload.Status,
		UserNumber:   user.Number,
	}

	var duplicateFeedingGroup models.FeedingGroup
	if err := config.DB.Where("name = ?", payload.Name).First(&duplicateFeedingGroup).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Feeding group name already exists"})
		return
	}

	if err := config.DB.Create(&saveFeedingGroup).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save feeding group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feeding group has been created"})
}

func FetchFeedingGroups(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var feedingGroups []models.FeedingGroup
	if err := config.DB.Preload("Campus").Preload("User").Find(&feedingGroups).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch feeding groups"})
		return
	}

	var response []responses.FeedingGroupResponse
	for _, feedingGroup := range feedingGroups {
		response = append(response, responses.FeedingGroupResponse{
			Number:       feedingGroup.Number,
			Name:         feedingGroup.Name,
			CampusNumber: feedingGroup.Campus.Number,
			CampusName:   feedingGroup.Campus.Name,
			Photo:        feedingGroup.User.Photo,
			UserName:     feedingGroup.User.FirstName + feedingGroup.User.SecondName + feedingGroup.User.LastName,
			Status:       feedingGroup.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
