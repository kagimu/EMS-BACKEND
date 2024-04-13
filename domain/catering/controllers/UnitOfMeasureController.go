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

func CreateUnitOfMeasure(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.UnitOfMeasureRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveUnitOfMeaure := models.UnitOfMeasure{
		Number:     uuid.NewV4().String(),
		Name:       payload.Name,
		Code:       payload.Code,
		UserNumber: user.Number,
	}

	var duplicate models.UnitOfMeasure
	if err := config.DB.Where("name = ?", payload.Name).First(&duplicate).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unit of measure with this name already exists"})
		return
	}

	if err := config.DB.Where("code = ?", payload.Code).First(&duplicate).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unit of measure with this code already exists"})
		return
	}

	if err := config.DB.Create(&saveUnitOfMeaure).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create unit of measure"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unit of measure created successfully"})
}

func FetchUnitsOfMeasure(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var unitsOfMeasure []models.UnitOfMeasure
	if err := config.DB.Find(&unitsOfMeasure).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch units of measure"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Units of measure fetched successfully", "data": unitsOfMeasure})
}
