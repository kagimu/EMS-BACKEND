package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/api/requests"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/services"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

var itemService = services.NewItemService()

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

	err = itemService.SaveService(user, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save service"})
		return
	}

	/*
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
			saveCampusService := models2.CampusItem{
				Number:       uuid.NewV4().String(),
				CampusNumber: campusNumber,
				ItemNumber:   saveService.ItemNumber,
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
		}*/

	c.JSON(http.StatusCreated, gin.H{"message": "Service created successfully"})
}

func FetchServices(c *gin.Context) {

	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	schoolNumber := c.Query("school_number")

	if user.Type != constants.TypeClient {
		if schoolNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Please select the school you want to view services for"})
			return
		}
	}

	response, err := itemService.AllServiceItems(schoolNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func FetchAllItems(c *gin.Context) {

	schoolNumber := c.Query("school_number")
	if schoolNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Provide school number"})
		return
	}

	itemType := c.Query("type")

	response, err := itemService.AllItems(schoolNumber, itemType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func EditService(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if user.Number == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized!"})
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
	err = itemService.UpdateService(serviceNumber, user, payload)

	c.JSON(http.StatusOK, gin.H{"message": "Service updated successfully"})
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
	if err := tx.Where("item_number = ?", serviceNumber).Delete(&models2.CampusItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete campus services", "error": err.Error()})
		return

	}

	if err := tx.Delete(&service).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete service", "error": err.Error()})
		return
	}
	if err := tx.Where("number = ?", serviceNumber).Delete(&models.Item{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete service item", "error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to commit service transaction", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted"})
}
