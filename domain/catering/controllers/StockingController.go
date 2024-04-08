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

func CreateStocking(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.StockingRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}

	for _, stock := range payload.ItemPriceAndQuantity {
		saveStock := models.Stocking{
			Number:         uuid.NewV4().String(),
			CampusNumber:   payload.CampusNumber,
			UnitCost:       stock.UnitCost,
			MenuItemNumber: stock.MenuItemNumber,
			Quantity:       stock.Quantity,
			UserNumber:     user.Number,
		}

		if err := config.DB.Create(&saveStock).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to create interview grade"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stocking created successfully"})
}

func FetchStockings(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var stockings []models.Stocking
	if err := config.DB.Find(&stockings).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch item assemblies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stockings})
}
