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

func CreateFeedingSessionItemController(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.FeedingSessionItemRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}

	for _, feedingItem := range payload.MenuItemQuantity {
		saveItemAssembly := models.FeedingSessionItem{
			Number:               uuid.NewV4().String(),
			FeedingSessionNumber: payload.FeedingSessionNumber,
			MenuItemNumber:       feedingItem.MenuItemNumber,
			Quantity:             feedingItem.Quantity,
			UserNumber:           user.Number,
		}

		if err := config.DB.Create(&saveItemAssembly).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create feeding session item"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feeding session item created successfully"})
}

func FetchFeedingSessionItems(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var feedingSessionItems []models.FeedingSessionItem
	if err := config.DB.Preload("FeedingSession").Preload("MenuItem").Find(&feedingSessionItems).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch item assemblies"})
		return
	}

	var response []responses.FeedingSessionItemResponse
	for _, feedingSessionItem := range feedingSessionItems {
		response = append(response, responses.FeedingSessionItemResponse{
			Number:               feedingSessionItem.Number,
			FeedingSessionNumber: feedingSessionItem.FeedingSession.Number,
			FeedingSessionName:   feedingSessionItem.FeedingSession.Name,
			FeedingSessionStatus: feedingSessionItem.FeedingSession.Status,
			MenuItemNumber:       feedingSessionItem.MenuItemNumber,
			MenuItemName:         feedingSessionItem.MenuItem.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
