package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/api/requests"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/models"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateDiscount(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ServiceDiscountRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input- " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var service models2.Service
	if err := config.DB.Where("number = ?", payload.ServiceNumber).First(&service).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Service not found"})
		return
	}

	var authorizor authModel.User
	if err := config.DB.Where("number = ?", payload.AuthorizedBy).First(&authorizor).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var student studentModel.Student
	if err := config.DB.Where("number = ?", payload.StudentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	tx := config.DB.Begin()

	saveServiceDiscount := models2.Discount{
		Number:         uuid.NewV4().String(),
		ServiceNumber:  payload.ServiceNumber,
		AuthorizedBy:   payload.AuthorizedBy,
		StudentNumber:  payload.StudentNumber,
		DiscountAmount: payload.DiscountAmount,
		Reason:         payload.Reason,
		UserNumber:     user.Number,
	}

	if err := tx.Create(&saveServiceDiscount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create service discount"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save service discount"})
		return
	}

	c.JSON(http.StatusCreated, saveServiceDiscount)

}
