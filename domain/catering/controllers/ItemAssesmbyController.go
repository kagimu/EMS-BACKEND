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

func CreateItemAssembly(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.ItemAssesmbly
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	for i, primaryItem := range payload.PrimaryItemNumber {

		saveItemAssembly := models.ItemAssesmbly{
			Number:              uuid.NewV4().String(),
			PrimaryItemNumber:   primaryItem,
			SecondaryItemNumber: payload.SecondaryItemNumber[i],
			Quantity:            payload.Quantity,
			UserNumber:          user.Number,
		}

		if err := config.DB.Create(&saveItemAssembly).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to create interview grade"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item assembly created successfully"})
}

func FetchItemAssemblies(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var itemAssesmblies []models.ItemAssesmbly
	if err := config.DB.Find(&itemAssesmblies).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch item assemblies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": itemAssesmblies})
}
