package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
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

	if err = payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.ProcessNewStock(payload, user.Number)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Stocking error: " + err.Error()})
		return
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
	if err := config.DB.Preload("MenuItem").Preload("Campus").Find(&stockings).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch item assemblies"})
		return
	}

	var response []responses.StockingResponse
	for _, stocking := range stockings {
		response = append(response, responses.StockingResponse{
			Number:         stocking.Number,
			Quantity:       stocking.Quantity,
			UnitCost:       stocking.UnitCost,
			MenuItemNumber: stocking.MenuItemNumber,
			CampusNumber:   stocking.CampusNumber,
			CampusName:     stocking.Campus.Name,
			MenuItemName:   stocking.MenuItem.Name,
			MenuItemType:   stocking.MenuItem.Type,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func FetchItemStock(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	stockItemNumber := c.Param("stock_item_number")
	if stockItemNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Stock item number is required"})
		return
	}

	var stocking models.Stocking
	if err := config.DB.Where("number = ?", stockItemNumber).First(&stocking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Stocking not found"})
		return
	}

	response := responses.StockingResponse{
		Number:         stocking.Number,
		Quantity:       stocking.Quantity,
		UnitCost:       stocking.UnitCost,
		MenuItemNumber: stocking.MenuItemNumber,
		CampusNumber:   stocking.CampusNumber,
		CampusName:     stocking.Campus.Name,
		MenuItemName:   stocking.MenuItem.Name,
		MenuItemType:   stocking.MenuItem.Type,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
