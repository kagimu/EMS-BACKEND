package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	helpers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/gin-gonic/gin"
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
		Preload("Parent").
		Preload("PreviousSchools").
		Where("number = ?", applicationNumber).
		First(&application).
		Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     err.Error(),
			"message": "Application not found",
		})
		return
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
	query := config.DB.Model(&models.Application{}).Preload("Student").Preload("Parent").Preload("PreviousSchools").Find(&applications)

	_, err := helpers.Paginate(query, paginationParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Pagination failed",
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
			"message": "application fetching failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    applications,
		"message": "applications fetched successfully",
	})
}
