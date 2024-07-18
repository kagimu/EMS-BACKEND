package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
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

	unitOfMeasureService := services.NewUnitOfMeasureService(user)
	err = unitOfMeasureService.NewUnitOfMeasure(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unit of measure created successfully"})
}

func FetchUnitsOfMeasure(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	schoolNumber := c.Query("school_number")
	if schoolNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Choose the school"})
		return
	}

	unitOfMeasureService := services.NewUnitOfMeasureService(user)
	list, err := unitOfMeasureService.FindUnitOfMeasures(schoolNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Units of measure fetched successfully", "data": list})
}

func EditUnitOfMeasure(c *gin.Context) {

	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.EditUnitOfMeasureRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payload error: " + err.Error()})
		return
	}

	unitOfMeasureService := services.NewUnitOfMeasureService(user)
	err = unitOfMeasureService.UpdateUnitOfMeasure(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unit of measure updated successfully"})

}

func RemoveUnitOfMeasure(c *gin.Context) {

	unitOfMeasureNumber := c.Param("unit_of_measure_number")
	if unitOfMeasureNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Choose unit of measure to remove"})
		return
	}

	unitOfMeasureService := services.NewUnitOfMeasureService(nil)
	err := unitOfMeasureService.RemoveUnitOfMeasures(unitOfMeasureNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unit of measure removed successfully"})
}
