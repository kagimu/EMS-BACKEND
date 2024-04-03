package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	helpers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FetchApplication(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	var application = models.Application{}
	var applicationNumber = c.Param("application_number")

	if err := config.DB.Model(&models.Application{}).
		Preload("Student").
		Preload("Student.CampusClass").
		Preload("Parent").
		Preload("PreviousSchools").
		Where("number = ?", applicationNumber).
		First(&application).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Application not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, application)
}

func FetchApplications(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	paginationParams := helpers.ParsePaginationParams(c)

	var applications []models.Application
	query := config.DB.Model(&models.Application{}).Preload("Student").Preload("Parent").Preload("PreviousSchools")

	_, err := helpers.Paginate(query, paginationParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Pagination failed",
		})
		return
	}

	if err := query.Find(&applications).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "application fetching failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"data":    applications,
		"message": "applications fetched successfully",
	})
}

func FetchApplicationsByStatus(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	status := c.Query("status")

	var applications []models.Application
	if err := config.DB.Where("status = ?", status).Find(&applications).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Application fetching failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    applications,
		"message": "Applications fetched successfully",
	})
}

func EditApplication(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	type ApplicationInfo struct {
		Status                string `json:"status"`
		ApplicationFeePayment bool   `json:"application_fee_payment"`
		InterviewMode         string `json:"interview_mode"`
	}

	var payload ApplicationInfo
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if payload.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Status is required",
		})
		return
	}

	var application models.Application
	var applicationNumber = c.Param("application_number")

	if err := config.DB.First(&application, "number = ?", applicationNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No application found with the specified number"})
		return
	}

	updateApplication := models.Application{
		Status:                payload.Status,
		ApplicationFeePayment: payload.ApplicationFeePayment,
		InterviewMode:         payload.InterviewMode,
	}

	if err := config.DB.Model(&application).Where("number = ?", applicationNumber).Updates(updateApplication).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    application,
		"message": "Application updated successfully",
	})
}

func FetchApplicationByNumber(c *gin.Context, applicationNumber string) (*models.Application, error) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return nil, err
	}

	var application models.Application
	if err := config.DB.Where("number = ?", applicationNumber).Preload("Parent").First(&application).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Application not found"})
		return nil, err
	}

	return &application, nil
}
