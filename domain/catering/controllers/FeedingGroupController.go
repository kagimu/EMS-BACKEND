package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
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

	saveFeedingGroup := models.FeedingGroup{
		Number:       uuid.NewV4().String(),
		Name:         payload.Name,
		SchoolNumber: payload.SchoolNumber,
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

	schoolNumber := c.Query("school_number")
	if schoolNumber == "" {
		c.JSON(http.StatusBadRequest, commons.Feedback("Provide school number"))
		return
	}

	status := c.Query("status")

	feedingGroups, err := services.GetFeedingGroupsBySchoolNumber(schoolNumber, status)

	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Getting feeding groups error: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": feedingGroups})
}

func UpdateFeedingGroup(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.UpdateFeedingGroupRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Check the provided body request."})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.UpdateFeedingGroup(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Updating feeding group error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feeding group has been updated"})
}
