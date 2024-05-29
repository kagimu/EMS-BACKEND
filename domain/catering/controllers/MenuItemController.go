package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
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

	err = services.SaveMenuItem(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Menu item saved successfully"})
}

func FetchMenuItems(c *gin.Context) {
	var err error

	campusNumber := c.Query("campus_number")
	schoolNumber := c.Query("school_number")

	if campusNumber == "" && schoolNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Provide at least a campus or school number"})
		return
	}

	paginatedResponse, err := services.GetAllMenuItems(schoolNumber, campusNumber, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch menu items: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, paginatedResponse)
}

func EditMenuItem(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.ItemAssembly
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	menuItemNumber := c.Param("menu_item_number")
	var menuItem models.MenuItem
	if err := config.DB.Where("number = ?", menuItemNumber).First(&menuItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Menu item not found"})
		return
	}

	updateMenuItem := models.MenuItem{
		Name:                payload.Name,
		Code:                payload.Code,
		Type:                payload.Type,
		UnitOfMeasureNumber: payload.UnitOfMeasureNumber,
	}

	if err := config.DB.Model(&menuItem).Where("number = ?", menuItemNumber).Updates(updateMenuItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update menu item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": menuItem})
}

func FetchMenuItemDetails(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	itemNumber := c.Param("menu_item_number")
	var menuItem models.MenuItem
	if err := config.DB.Where("number = ?", itemNumber).First(&menuItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Menu item not found"})
		return
	}

	response, err := services.MenuItemDetails(itemNumber)

	c.JSON(http.StatusOK, gin.H{"data": response})
}
