package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
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
	parentNumber := c.Param("parent_number")
	studentNumber := c.Param("student_number")

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
		StudentNumber:         studentNumber,
		ParentNumber:          parentNumber,
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

	applicationNumber := c.Param("application_number")

	var application models.Application

	err := config.DB.
		Where("number = ?", applicationNumber).
		Preload("Parent").
		Preload("Student").
		Preload("Student.Class").
		Preload("Student.PreferredCampus").
		Preload("Student.SourceOfSchoolInfo").
		First(&application).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     err.Error(),
			"message": "Application not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "application fetched successfully",
		"data":    application,
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
			"message": "Applications not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    applications,
		"message": "applications fetched successfully",
	})
}
