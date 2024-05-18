package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/api/requests"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/models"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
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

	generateProductSerialCode := models2.SerialCode{}

	if err := config.DB.Create(&generateProductSerialCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate service serial code"})
	}

	payload.SerialCode = fmt.Sprintf("SVC%04d", generateProductSerialCode.ID)

	tx := config.DB.Begin()
	itemNumber := uuid.NewV4().String()
	saveItem := models.Item{
		Number:          itemNumber,
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

	if err := tx.Create(&saveItem).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create item " + err.Error()})
		return
	}

	saveService := models2.Service{
		ItemNumber: itemNumber,
	}

	if err := tx.Create(&saveService).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create service " + err.Error()})
		return
	}

	for _, campusNumber := range payload.CampusNumbers {
		saveCampusService := models2.CampusServices{
			Number:        uuid.NewV4().String(),
			CampusNumber:  campusNumber,
			ServiceNumber: saveService.ItemNumber,
		}

		if err := tx.Create(&saveCampusService).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create service for the campus"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to commit service transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Service created successfully"})
}

func FetchServices(c *gin.Context) {

	var services []models2.Service
	if err := config.DB.Preload("Item").Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch services"})
		return
	}

	c.JSON(http.StatusOK, services)
}

func FetchAllItems(c *gin.Context) {
	var items = make([]models.Item, 0)

	var item []models.Item
	if err := config.DB.Find(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch items"})
		return
	}

	for _, service := range item {
		items = append(items, service)
	}
	c.JSON(http.StatusOK, items)
	return
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

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input %s:", err.Error())})
		return

	}

	serviceNumber := c.Param("service_number")
	var service models.Item
	if err := config.DB.Where("number = ?", serviceNumber).First(&service).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Service not found"})
		return
	}

	updateService := models.Item{
		Description:     payload.Description,
		Cost:            payload.Cost,
		Name:            payload.Name,
		Mandatory:       payload.Mandatory,
		SerialCode:      payload.SerialCode,
		RepeatCycle:     payload.RepeatCycle,
		PaymentMode:     payload.PaymentMode,
		OneTimePurchase: payload.OneTimePurchase,
	}

	if err := config.DB.Model(&service).Where("number = ?", serviceNumber).Updates(updateService).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update service item"})
		return
	}

	// synchronize campus services by deleting all and adding them again
	if err := config.DB.Where("service_number = ?", serviceNumber).Delete(&models2.CampusServices{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete campus services"})
		return

	}

	for _, campusNumber := range payload.CampusNumbers {
		saveCampusService := models2.CampusServices{
			Number:        uuid.NewV4().String(),
			CampusNumber:  campusNumber,
			ServiceNumber: serviceNumber,
		}
		if err := config.DB.Create(&saveCampusService).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create campus service"})
			return
		}

	}

	c.JSON(http.StatusOK, service)
}

func DeleteService(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	serviceNumber := c.Param("service_number")
	var service models2.Service
	if err := config.DB.Where("item_number = ?", serviceNumber).First(&service).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Service not found"})
		return
	}

	//start db transaction
	tx := config.DB.Begin()

	//delete campus services
	if err := tx.Where("service_number = ?", serviceNumber).Delete(&models2.CampusServices{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete campus services"})
		return

	}

	if err := tx.Delete(&service).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete service"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to commit service transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted"})
}
