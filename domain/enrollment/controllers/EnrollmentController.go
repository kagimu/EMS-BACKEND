package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateEnrollment(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.EnrollmentRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveEnrollment := models.Enrollment{
		Number:            uuid.NewV4().String(),
		StudentNumber:     payload.StudentNumber,
		CampusNumber:      payload.CampusNumber,
		SchoolNumber:      payload.SchoolNumber,
		EnrollmentDate:    payload.EnrollmentDate,
		StudyPeriodNumber: payload.StudyPeriodNumber,
		Status:            payload.Status,
		UserNumber:        user.Number,
	}

	if err := config.DB.Create(&saveEnrollment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save enrollment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Enrollment created successfully"})
}

func FetchEnrolledStudents(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var enrolledStudents []models.Enrollment
	if err := config.DB.Preload("School").Preload("Campus").Preload("Student").Preload("StudyPeriod").Find(&enrolledStudents).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch enrolled students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": enrolledStudents})
}

func EnrolledStudentsByCampus(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var enrolledStudents []models.Enrollment
	campusNumber := c.Query("campus_number")

	if err := config.DB.Where("campus_number = ?", campusNumber).Find(&enrolledStudents).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch enrolled students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": enrolledStudents})
}
