package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	helpers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/preAdmissionInterviewTest/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/preAdmissionInterviewTest/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreatePreAdmissionInterviewReport(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var payload requests.PreAdmissionInterviewTestRequest

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
		Number:                          uuid.NewV4().String(),
		ApplicationNumber:               payload.ApplicationNumber,
		AdmissionStatus:                 payload.AdmissionStatus,
		ClassNumber:                     payload.ClassNumber,
		AttentionNeeded:                 payload.AttentionNeeded,
		ReasonForChangingAdmissionClass: payload.ReasonForChangingAdmissionClass,
		ReasonForNotAdmittingStudent:    payload.ReasonForNotAdmittingStudent,
		UserNumber:                      user.Number,
	}

	if err := config.DB.Create(&preAdmissionInterviewReport).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create pre admission interview report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pre admission interview report saved successfully", "data": preAdmissionInterviewReport})
}

func GetAllPreAdmissionInterviewReports(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	paginationParams := helpers.ParsePaginationParams(c)

	var preAdmissionInterviewReports []models.PreAdmissionInterviewTestReport
	query := config.DB.Model(&models.PreAdmissionInterviewTestReport{})

	_, err := helpers.Paginate(query, paginationParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Pagination failed",
		})
		return
	}

	if err := query.Find(&preAdmissionInterviewReports).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "pre admission interview test report  fetching failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pre admission interview reports fetched successfully", "data": preAdmissionInterviewReports})
}

func FetchPreAdmissionInterviewReportByNumber(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found"})
		return
	}

	var preAdmissionInterviewReport models.PreAdmissionInterviewTestReport
	preAdmissionInterviewReportNumber := c.Param("number")

	if preAdmissionInterviewReportNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Pre admission interview report number is required",
		})
		return
	}

	if err := config.DB.First(&preAdmissionInterviewReport, "number = ?", preAdmissionInterviewReportNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "No pre admission interview report found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": preAdmissionInterviewReport})

}
