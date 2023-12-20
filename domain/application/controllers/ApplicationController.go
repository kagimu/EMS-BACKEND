package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	parentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	previousSchoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchool/models"
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

	applicationNumber := c.Param("application_number")
	var application models.Application

	if err := config.DB.Where("number = ?", applicationNumber).First(&application).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "application not found"})
		return
	}

	// TODO: Improve this
	var parents parentModels.Parent
	config.DB.Raw(`
	SELECT * FROM parents WHERE parents.application_number = ?`, applicationNumber).Scan(&parents)

	var student studentModel.Student
	err := config.DB.Raw(`
	SELECT * FROM previous_schools WHERE previous_schools.application_number = ?`, applicationNumber).Scan(&student).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "student   not found"})
		return
	}

	var previousSchool previousSchoolModel.PreviousSchool
	config.DB.Raw(`
	SELECT * FROM previous_schools WHERE previous_schools.application_number = ?`, applicationNumber).Scan(&previousSchool)

	var previousSchoolClasses []previousSchoolModel.PreviousSchoolClass
	err = config.DB.Raw(`
	SELECT * FROM previous_school_classes WHERE previous_school_classes.student_number = ?`, student.Number).Preload("PreviousSchoolClass").Scan(&previousSchoolClasses).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Err":     student.Number,
			"message": "application   not found"})
		return
	}
	responseData := gin.H{
		"message": "application fetched successfully",
		"data": gin.H{
			"parents":                 parents,
			"application":             application,
			"student":                 student,
			"previous_school_classes": previousSchoolClasses,
			"previous_school":         previousSchool,
		},
	}

	c.JSON(http.StatusOK, responseData)
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
