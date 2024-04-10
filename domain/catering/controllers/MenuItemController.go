package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
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

	if payload.Type != constants.TYPE_ASSEMBLED && payload.Type != constants.TYPE_ITEM {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid menu item type"})
		return
	}

	saveMenuItem := models.MenuItem{
		Number:              uuid.NewV4().String(),
		Name:                payload.Name,
		Code:                payload.Code,
		Type:                payload.Type,
		CampusNumber:        payload.CampusNumber,
		UnitOfMeasureNumber: payload.UnitOfMeasureNumber,
		UserNumber:          user.Number,
	}

	if err := config.DB.Create(&saveMenuItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create menu item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Menu item saved successfully"})
}

func FetchMenuItems(c *gin.Context) {
	var response []responses.MenuItemResponse

	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var menuItems []models.MenuItem
	if err := config.DB.Preload("Campus").Preload("UnitOfMeasure").Find(&menuItems).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch menu items"})
		return
	}

	for _, menuItem := range menuItems {
		response = append(response, responses.MenuItemResponse{
			Number:              menuItem.Number,
			Name:                menuItem.Name,
			Code:                menuItem.Code,
			Type:                menuItem.Type,
			Email:               menuItem.Campus.Email,
			CampusName:          menuItem.Campus.Name,
			CampusNumber:        menuItem.CampusNumber,
			Location:            menuItem.Campus.Location,
			UnitOfMeasureNumber: menuItem.UnitOfMeasureNumber,
			UnitOfMeasureName:   menuItem.UnitOfMeasure.Name,
			UnitOfMeasureCode:   menuItem.UnitOfMeasure.Code,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
