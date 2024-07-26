package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	campusNumber := c.Param("campus_number")
	itemNumber := c.Param("item_number")

	if campusNumber == "" && itemNumber == "" {
		c.JSON(http.StatusBadRequest, commons.Feedback("Choose campus or menu item"))
		return
	}

	stock, err := services.AllStocks(campusNumber, itemNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stock})
}

func FetchItemStock(c *gin.Context) {

	stockItemNumber := c.Param("stock_item_number")
	if stockItemNumber == "" {
		c.JSON(http.StatusBadRequest, commons.Feedback("Stock item number is required"))
		return
	}

	stock, err := services.GetStockItem(stockItemNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stock})
}
