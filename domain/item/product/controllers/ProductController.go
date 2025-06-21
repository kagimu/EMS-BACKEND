package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/services"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/product/api/requests"
	productModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/product/models"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

var itemService = services.NewItemProduct()

func SaveProduct(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ProductRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf(" %s", err.Error())})
		return
	}

	err = itemService.SaveProduct(user, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save product: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added successfully"})
}

func EditProduct(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ProductRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf(" %s", err.Error())})
		return
	}

	productNumber := c.Param("product_number")
	if productNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Product number is required"})
		return
	}

	err = itemService.UpdateProduct(productNumber, user, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save product: " + err.Error()})
		return
	}

	/*
		var product models.Item
		if err := config.DB.First(&product, "number = ?", productNumber).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found."})
			return
		}

		checkSerial, err := repositories.GetItemBySerialCode(payload.SerialCode, payload.SchoolNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if checkSerial.Number != "" && checkSerial.Number != productNumber {
			c.JSON(http.StatusBadRequest, gin.H{"message": "This serial code is already in use."})
			return
		}

		txn := config.DB.Begin()

		updateProduct := models.Item{
			Name:         payload.Name,
			Cost:         payload.Cost,
			Description:  payload.Description,
			UserNumber:   user.Number,
			SerialCode:   payload.SerialCode,
			Mandatory:    payload.Mandatory,
			SchoolNumber: payload.SchoolNumber,
		}

		if err := txn.Model(&product).Where("number = ?", productNumber).Updates(updateProduct).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update product"})
			return
		}

		// synchronize campus services by deleting all and adding them again
		if err := txn.Where("item_number = ?", productNumber).Delete(&models2.CampusItem{}).Error; err != nil {
			txn.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete campus services"})
			return
		}

		for _, campusNumber := range payload.CampusNumbers {
			saveCampusService := models2.CampusItem{
				Number:       uuid.NewV4().String(),
				CampusNumber: campusNumber,
				ItemNumber:   productNumber,
			}
			if err := txn.Create(&saveCampusService).Error; err != nil {
				txn.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create campus service"})
				return
			}
		}

		txn.Commit()*/

	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

func FetchProducts(c *gin.Context) {

	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	//var response []responses.ServiceResponse

	schoolNumber := c.Query("school_number")

	if user.Type != constants.TypeClient {
		if schoolNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Please select the school you want to view services for"})
			return
		}
	}

	response, err := itemService.AllProductItems(schoolNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func DeleteProduct(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	productNumber := c.Param("product_number")
	var service productModels.Product
	if err := config.DB.Where("number = ?", productNumber).First(&service).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found "})
		return
	}

	if err := config.DB.Delete(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
