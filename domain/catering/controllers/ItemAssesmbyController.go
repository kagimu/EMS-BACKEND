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

func CreateItemAssembly(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.ItemAssembly
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tx := config.DB.Begin()

	menuItemNumber := uuid.NewV4().String()
	saveMenuItem := models.MenuItem{
		Number:              menuItemNumber,
		Name:                payload.Name,
		Code:                payload.Code,
		Type:                payload.Type,
		CampusNumber:        payload.CampusNumber,
		UnitOfMeasureNumber: payload.UnitOfMeasureNumber,
	}

	if err := tx.Create(&saveMenuItem).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create menu item"})
		return
	}

	for _, secondaryItem := range payload.SecondaryItemQuantity {
		saveItemAssembly := models.ItemAssembly{
			Number:              uuid.NewV4().String(),
			PrimaryItemNumber:   menuItemNumber,
			SecondaryItemNumber: secondaryItem.SecondaryItemNumber,
			Quantity:            secondaryItem.Quantity,
			UserNumber:          user.Number,
		}

		if err := tx.Create(&saveItemAssembly).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create interview grade"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create assesmbly items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item assembly created successfully"})
}

func FetchItemAssemblies(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var itemAssesmblies []models.ItemAssembly
	if err := config.DB.Preload("PrimaryItem").Preload("SecondaryItem").Find(&itemAssesmblies).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch item assemblies"})
		return
	}

	var response []responses.ItemAssemblyResponse
	for _, itemAssembly := range itemAssesmblies {
		response = append(response, responses.ItemAssemblyResponse{
			Number:              itemAssembly.Number,
			PrimaryItemNumber:   itemAssembly.PrimaryItemNumber,
			PrimaryItemName:     itemAssembly.SecondaryItem.Name,
			SecondaryItemNumber: itemAssembly.SecondaryItemNumber,
			SecondaryItemName:   itemAssembly.SecondaryItem.Name,
			Quantity:            itemAssembly.Quantity,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
