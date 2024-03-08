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

func CreateTaskSupervision(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.TaskSupervisionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	saveTaskSupervision := models.TaskSupervision{
		Number:              uuid.NewV4().String(),
		DateOfSupervision:   payload.DateOfSupervision,
		DescriptiveAnalysis: payload.DescriptiveAnalysis,
		TaskNumber:          payload.TaskNumber,
		UserNumber:          user.Number,
	}

	if err := config.DB.Create(&saveTaskSupervision).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save task supervision"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task supervision created successfully"})
}

func FetchTaskSUpervisions(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var taskSupervisions []models.TaskSupervision
	if err := config.DB.Where("user_number = ?", user.Number).Find(&taskSupervisions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch task supervisions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task supervisions fetched successfully", "data": taskSupervisions})
}
