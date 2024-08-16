package controllers

import (
	"fmt"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/api/responses"
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
		SchoolNumber:    payload.SchoolNumber,
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
	var response []responses.ServiceResponse
	for _, service := range services {
		var campusServices []models2.CampusServices
		if err := config.DB.Where("service_number = ?", service.ItemNumber).Find(&campusServices).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch campus services"})
			return
		}

		var campusServiceResponse []responses.CampusServiceResponse
		for _, campusService := range campusServices {
			//get campus name
			var campus campusModel.Campus
			if err := config.DB.Where("number = ?", campusService.CampusNumber).First(&campus).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch campus"})
				return
			}

			campusServiceResponse = append(campusServiceResponse, responses.CampusServiceResponse{
				Number:       campusService.Number,
				CampusNumber: campusService.CampusNumber,
				Name:         campus.Name,
			})
		}

		response = append(response, responses.ServiceResponse{
			ItemNumber: service.ItemNumber,
			Item: responses.ItemResponse{
				Number:          service.Item.Number,
				Name:            service.Item.Name,
				Description:     service.Item.Description,
				Type:            service.Item.Type,
				Cost:            service.Item.Cost,
				SerialCode:      service.Item.SerialCode,
				PaymentMode:     service.Item.PaymentMode,
				Mandatory:       service.Item.Mandatory,
				SchoolNumber:    service.Item.SchoolNumber,
				UserNumber:      service.Item.UserNumber,
				RepeatCycle:     service.Item.RepeatCycle,
				OneTimePurchase: service.Item.OneTimePurchase,
			},
			Campuses: campusServiceResponse,
		})
	}

	c.JSON(http.StatusOK, response)
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
		SchoolNumber:    payload.SchoolNumber,
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
	if err := tx.Where("number = ?", serviceNumber).Delete(&models.Item{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete service item"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to commit service transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted"})
}
