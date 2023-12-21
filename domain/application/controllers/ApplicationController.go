package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	parentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateApplication(c *gin.Context) {

	type ApplicationRequest struct {
		Status                string `json:"status"`
		ApplicationFeePayment bool   `json:"application_fee_payment"`
		StudentNumber         string `json:"student_number"`
		ParentNumber          string `json:"parent_number"`
	}

	var application ApplicationRequest

	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal server error"})
		return
	}

	if application.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "status is required"})
		return
	}

	userNumber := c.Request.Header["User-Number"][0]

	var user authModel.User
	if err := config.DB.
		Where("number = ?", userNumber).
		First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	tx := config.DB.Begin()
	saveApplication := models.Application{
		Number:                uuid.NewV4().String(),
		Status:                string(application.Status),
		UserNumber:            userNumber,
		ApplicationFeePayment: application.ApplicationFeePayment,
	}

	if err := tx.Create(&saveApplication).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to save application",
			"details": err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to save application",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "application created successfully",
		"data":    saveApplication,
	})
}

func FetchApplication(c *gin.Context) {

	userNumber := c.Request.Header["User-Number"][0]

	var user authModel.User
	if err := config.DB.
		Where("number = ?", userNumber).
		First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     err.Error(),
			"message": "User not found"})
		return
	}

	var application models.Application
	var applicationNumber = c.Param("application_number")
	var student studentModel.Student
	var parent parentModel.Parent

	if err := config.DB.Table("applications").
		Select("applications.*, students.*, parents.*").
		Joins("RIGHT JOIN students ON students.application_number = applications.number").
		Joins("RIGHT JOIN parents ON parents.application_number = applications.number").
		Where("applications.number = ?", applicationNumber).
		Scan(&student).
		Scan(&parent).
		Scan(&application).
		Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     err.Error(),
			"message": "Application not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"application": application,
			"student":     student,
			"parent":      parent,
		},
	})

}

func FetchApplications(c *gin.Context) {
	userNumber := c.Request.Header["User-Number"][0]
	var user authModel.User

	if err := config.DB.
		Where("number = ?", userNumber).
		First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     err.Error(),
			"message": "User not found"})
		return
	}

	var applications []models.Application
	if err := config.DB.Find(&applications).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     err.Error(),
			"message": "Application fetching failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    applications,
		"message": "applications fetched successfully",
	})
}
