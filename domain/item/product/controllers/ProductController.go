package controllers

import (
	"fmt"
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

	generateProductSerialCode := productModels.ProductSerialCode{}

	if err := tx.Create(&generateProductSerialCode).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": " Failed to generate product serial code"})
		return
	}

	payload.SerialCode = fmt.Sprintf("PDT%04d", generateProductSerialCode.ID)

	itemNumber := uuid.NewV4().String()
	saveProductItem := models.Item{
		Number:             itemNumber,
		Name:               payload.Name,
		Cost:               payload.Cost,
		SerialCode:         payload.SerialCode,
		Description:        payload.Description,
		Prefix:             payload.Prefix,
		DiscountPercentage: payload.DiscountPercentage,
		DiscountAmount:     payload.DiscountAmount,
		UserNumber:         user.Number,
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

	updateProduct := models.Item{
		Number:             uuid.NewV4().String(),
		Name:               payload.Name,
		Cost:               payload.Cost,
		Description:        payload.Description,
		Prefix:             payload.Prefix,
		DiscountPercentage: payload.DiscountPercentage,
		DiscountAmount:     payload.DiscountAmount,
		UserNumber:         user.Number,
	}

	if err := config.DB.Model(&product).Where("number = ?", productNumber).Updates(updateProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, updateProduct)
}

func FetchProducts(c *gin.Context) {

	var products []productModels.Product
	if err := config.DB.Preload("Item").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func DeleteProduct(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	productNumber := c.Param("product_number")
	var service productModels.Product
	if err := config.DB.Where("number = ?", productNumber).First(&service).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	if err := config.DB.Delete(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
