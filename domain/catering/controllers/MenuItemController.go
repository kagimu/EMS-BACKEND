package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

func CreateMenuItem(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.MenuItemRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	_, err = services.SaveMenuItem(constants.TypeItem, payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Menu item saved successfully"})
}

func FetchMenuItemsPaginated(c *gin.Context) {
	var err error

	campusNumber := c.Query("campus_number")
	schoolNumber := c.Query("school_number")
	status := c.Query("status")
	itemType := c.Query("type")

	if campusNumber == "" && schoolNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Provide at least a campus or school number"})
		return
	}

	paginatedResponse, err := services.GetPaginatedMenuItems(schoolNumber, campusNumber, status, itemType, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch menu items: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, paginatedResponse)
}

func FetchMenuItems(c *gin.Context) {
	var err error

	campusNumber := c.Query("campus_number")
	schoolNumber := c.Query("school_number")
	status := c.Query("status")
	itemType := c.Query("type")

	if campusNumber == "" && schoolNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Provide at least a campus or school number"})
		return
	}

	response, err := services.GetAllMenuItems(schoolNumber, campusNumber, status, itemType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch menu items: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func EditMenuItem(c *gin.Context) {

	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.UpdateMenuItemRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.UpdateMenuItem(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Menu item updated successfully"})
}

func FetchMenuItemDetails(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	itemNumber := c.Query("menu_item_number")

	response, err := services.MenuItemDetails(itemNumber)

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func DeleteMenuItem(context *gin.Context) {

	_, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	itemNumber := context.Query("menu_item_number")

	err = services.DeleteMenuItem(itemNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Remove menu item error: " + err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": "Menu item removed successfully"})
}
