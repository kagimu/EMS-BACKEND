package controllers

import (
	"net/http"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	staffModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateFeeding(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.FeedingRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	staffNumber := payload.StaffNumber
	var staff staffModels.Staff
	if err := config.DB.First(&staff, "number = ?", staffNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No staff member found with the specified number"})
		return
	}

	layout := "2006-01-02T15:04:05.000Z"
	parsedDate, _ := time.Parse(layout, payload.Date)

	saveFeeding := models.Feeding{
		Number:               uuid.NewV4().String(),
		StaffNumber:          payload.StaffNumber,
		FeedingSessionNumber: payload.FeedingSessionNumber,
		StockingNumber:       payload.StockingNumber,
		Date:                 parsedDate,
		UserNumber:           user.Number,
	}

	if err := config.DB.Create(&saveFeeding).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save feeding"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feeding has been created"})
}

func FetchFeedings(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var feedings []models.Feeding
	if err := config.DB.Find(&feedings).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch feeding groups"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": feedings})
}
