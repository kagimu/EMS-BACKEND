package controllers

import (
	"fmt"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/models"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/product/api/requests"
	productModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/product/models"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateProduct(c *gin.Context) {
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

	tx := config.DB.Begin()

	itemNumber := uuid.NewV4().String()
	saveProductItem := models.Item{
		Number:      itemNumber,
		Name:        payload.Name,
		Cost:        payload.Cost,
		SerialCode:  payload.SerialCode,
		Description: payload.Description,
		UserNumber:  user.Number,
		Mandatory:   payload.Mandatory,
		Type:        "Product",
	}

	if err := tx.Create(&saveProductItem).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create product"})
		return
	}

	saveProduct := productModels.Product{
		ItemNumber: itemNumber,
	}

	if err := tx.Create(&saveProduct).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create product"})
		return
	}

	// loop through the campuses and create a campus service
	for _, campusNumber := range payload.CampusNumbers {
		campusService := models2.CampusService{
			ServiceNumber: itemNumber,
			CampusNumber:  campusNumber,
		}
		if err := tx.Create(&campusService).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create campus service"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, saveProduct)
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

	var product models.Item
	productNumber := c.Param("product_number")
	if err := config.DB.First(&product, "number = ?", productNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found."})
		return
	}
	txn := config.DB.Begin()

	updateProduct := models.Item{
		Number:      uuid.NewV4().String(),
		Name:        payload.Name,
		Cost:        payload.Cost,
		Description: payload.Description,
		UserNumber:  user.Number,
		SerialCode:  payload.SerialCode,
		Mandatory:   payload.Mandatory,
	}

	if err := txn.Model(&product).Where("number = ?", productNumber).Updates(updateProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product"})
		return
	}

	// synchronize campus services by deleting all and adding them again
	if err := txn.Where("service_number = ?", productNumber).Delete(&models2.CampusService{}).Error; err != nil {
		txn.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete campus services"})
		return
	}

	for _, campusNumber := range payload.CampusNumbers {
		saveCampusService := models2.CampusService{
			Number:        uuid.NewV4().String(),
			CampusNumber:  campusNumber,
			ServiceNumber: productNumber,
		}
		if err := txn.Create(&saveCampusService).Error; err != nil {
			txn.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create campus service"})
			return
		}
	}

	txn.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

func FetchProducts(c *gin.Context) {

	var products []productModels.Product
	if err := config.DB.Preload("Item").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch products"})
		return
	}
	var response []responses.ServiceResponse

	for _, product := range products {
		var campusServices []models2.CampusService
		if err := config.DB.Where("service_number = ?", product.ItemNumber).Find(&campusServices).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch campus services"})
			return
		}
		var campuses []responses.CampusServiceResponse
		for _, campusService := range campusServices {
			var campus campusModel.Campus
			if err := config.DB.First(&campus, "number = ?", campusService.CampusNumber).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch campus"})
				return
			}
			campuses = append(campuses, responses.CampusServiceResponse{
				Number:       campusService.Number,
				CampusNumber: campusService.CampusNumber,
				Name:         campus.Name,
			})
		}
		response = append(response, responses.ServiceResponse{
			ItemNumber: product.ItemNumber,
			Item: responses.ItemResponse{
				Number:          product.Item.Number,
				Name:            product.Item.Name,
				Description:     product.Item.Description,
				Type:            product.Item.Type,
				Cost:            product.Item.Cost,
				SerialCode:      product.Item.SerialCode,
				Prefix:          product.Item.Prefix,
				PaymentMode:     product.Item.PaymentMode,
				Mandatory:       product.Item.Mandatory,
				SchoolNumber:    product.Item.SchoolNumber,
				UserNumber:      product.Item.UserNumber,
				RepeatCycle:     product.Item.RepeatCycle,
				OneTimePurchase: product.Item.OneTimePurchase,
			},
			Campuses: campuses,
		})
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
