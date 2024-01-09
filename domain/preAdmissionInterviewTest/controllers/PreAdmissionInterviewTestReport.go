package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/preAdmissionInterviewTest/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreatePreAdmissionInterviewReport(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var payload models.PreAdmissionInterviewTestReport
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input payload"})
		return
	}

	if payload.ApplicationNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Application number is required"})
	}
	if payload.AdmissionStatus == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Admission status  is required"})
	}

	preAdmissionInterviewReport := models.PreAdmissionInterviewTestReport{
		Number:            uuid.NewV4().String(),
		ApplicationNumber: payload.ApplicationNumber,
		AdmissionStatus:   payload.AdmissionStatus,
		UserNumber:        user.Number,
	}

	if err := config.DB.Create(&preAdmissionInterviewReport).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create pre admission interview report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pre admission interview report saved successfully", "data": preAdmissionInterviewReport})

}
