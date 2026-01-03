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

func CreateTask(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.TaskRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input- " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	for _, person := range payload.PeopleInvolved {
		saveTask := models.Task{
			Number:         uuid.NewV4().String(),
			TimeFrame:      payload.TimeFrame,
			Title:          payload.Title,
			InputsUsed:     payload.InputsUsed,
			Description:    payload.Description,
			Outputs:        payload.Outputs,
			StaffNumber:    payload.StaffNumber,
			PeopleInvolved: person,
			UserNumber:     user.Number,
		}

		if err := config.DB.Create(&saveTask).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save task"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task created successfully"})
}

func EditTask(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.TaskRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input- " + err.Error()})
		return
	}

	var task models.Task
	taskNumber := c.Query("task_number")
	if err := config.DB.First(&task, "task_number = ?", task.Number).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No task found with the specified number"})
		return
	}

	saveTask := models.Task{
		TimeFrame:   payload.TimeFrame,
		Title:       payload.Title,
		InputsUsed:  payload.InputsUsed,
		Description: payload.Description,
		Outputs:     payload.Outputs,
		StaffNumber: payload.StaffNumber,
	}

	if err := config.DB.Model(&task).Where("number = ?", taskNumber).Updates(saveTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func DeleteTask(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	taskNumber := c.Param("task_number")
	var task models.Task

	if err := config.DB.First(&task, taskNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No task found"})
		return
	}

	if err := config.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "`Task` deleted successfully"})
}

func FetchTasks(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var tasks []models.Task
	if err := config.DB.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}
