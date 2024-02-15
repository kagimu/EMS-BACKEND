package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/service/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/service/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateService(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.ServiceRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input %s:", err.Error())})
		return
	}

	generateProductSerialCode := models.SerialCode{}

	if err := config.DB.Create(&generateProductSerialCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate service serial code"})
	}

	payload.SerialCode = fmt.Sprintf("SVC%04d", generateProductSerialCode.ID)

	tx := config.DB.Begin()

	saveService := models.Service{
		Number:          uuid.NewV4().String(),
		Cost:            payload.Cost,
		Description:     payload.Description,
		OneTimePurchase: payload.OneTimePurchase,
		UserNumber:      user.Number,
		SerialCode:      payload.SerialCode,
		Mandatory:       payload.Mandatory,
		RepeatCycle:     payload.RepeatCycle,
		PaymentMode:     payload.PaymentMode,
		Name:            payload.Name,
	}

	if err := tx.Create(&saveService).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create service " + err.Error()})
		return
	}

	for _, campusNumber := range payload.CampusNumbers {
		saveCampusService := models.CampusServices{
			Number:        uuid.NewV4().String(),
			CampusNumber:  campusNumber,
			ServiceNumber: saveService.Number,
		}

		if err := tx.Create(&saveCampusService).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create service for the campus"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save service"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Service created successfully"})
}

func FetchServices(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var services []models.Service
	if err := config.DB.Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch services"})
		return
	}

	c.JSON(http.StatusOK, services)
}

func EditService(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ServiceRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	serviceNumber := c.Param("service_number")
	var service models.Service
	if err := config.DB.Where("number = ?", serviceNumber).First(&service).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Service not found"})
		return
	}

	updateService := models.Service{
		Description: payload.Description,
		Cost:        payload.Cost,
		Name:        payload.Name,
	}

	if err := config.DB.Model(&service).Where("number = ?", serviceNumber).Updates(updateService).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update service"})
		return
	}

	c.JSON(http.StatusOK, service)
}

func DeleteService(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	serviceNumber := c.Param("service_number")
	var service models.Service
	if err := config.DB.Where("number = ?", serviceNumber).First(&service).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Service not found"})
		return
	}

	if err := config.DB.Delete(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted"})
}
