package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	helpers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	userfromheader "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/userFromHeader"
	parentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	previousModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchool/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
)

func FetchApplication(c *gin.Context) {
	_, err := userfromheader.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	var application models.Application
	var applicationNumber = c.Param("application_number")
	var student studentModel.Student
	var parent parentModel.Parent
	var previousSchools []previousModel.PreviousSchool

	if err := config.DB.Table("applications").
		Select("applications.*, students.*, parents.*, previous_schools.*").
		Joins("RIGHT JOIN students ON students.application_number = applications.number").
		Joins("RIGHT JOIN parents ON parents.application_number = students.application_number").
		Joins("RIGHT JOIN previous_schools ON previous_schools.application_number = applications.number").
		Where("applications.number = ?", applicationNumber).
		Scan(&student).
		First(&parent).
		First(&previousSchools).
		First(&application).
		Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     err.Error(),
			"message": "Application not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"application":     application,
			"student":         student,
			"parent":          parent,
			"previousSchools": previousSchools,
		},
	})
}

func FetchApplications(c *gin.Context) {
	if _, err := userfromheader.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	paginationParams := helpers.ParsePaginationParams(c)

	// Using the Student model and append the associated application
	// instead of the other way round

	var applications []models.Application
	query := config.DB.Model(&models.Application{}).Preload("Student").Preload("Parent").Preload("PreviousSchool").Find(&applications)

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
	_, err := userfromheader.GetUserFromHeader(c)
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
