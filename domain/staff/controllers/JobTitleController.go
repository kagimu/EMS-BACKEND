package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateJobTitle(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.JobTitleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveJobTitle := models.JobTitle{
		Number:      uuid.NewV4().String(),
		Name:        payload.Name,
		Description: payload.Description,
		UserNumber:  user.Number,
	}

	if err := config.DB.Create(&saveJobTitle).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save job title"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job title created successfully"})
}

func FetchJobTitles(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var jobTitles []models.JobTitle
	if err := config.DB.Find(&jobTitles).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch job titles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": jobTitles})
}

func EditJobTitle(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.JobTitleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	var jobTitle models.JobTitle
	jobTitleNumber := c.Query("job_title_number")
	if err := config.DB.First(&jobTitle, "number = ?", jobTitleNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No job title found with the specified number"})
		return
	}

	updateJobTitle := models.JobTitle{
		Name:        payload.Name,
		Description: payload.Description,
	}

	if err := config.DB.Model(&jobTitle).Where("number = ?", jobTitleNumber).Updates(updateJobTitle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update job title"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": jobTitle})
}
