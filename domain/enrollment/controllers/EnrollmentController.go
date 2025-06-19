package controllers

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/services"
	"gorm.io/gorm"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"github.com/gin-gonic/gin"
)

var enrolmentService = services2.NewEnrolmentService()

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

	err = enrolmentService.NewEnrolment(payload, user.Number)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Enrollment created successfully"})
}

func FetchEnrolledStudents(c *gin.Context) {

	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	paginatedResponse := pagination.PaginatedResponse{}

	schoolNumber := c.Query("school_number")
	streamNumber := c.Query("stream_number")
	campusNumber := c.Query("campus_number")
	campusClassNumber := c.Query("campus_class_number")
	studyPeriodNumber := c.Query("study_period_number")
	status := c.Query("status")

	if user.Type != constants.TypeClient {
		if schoolNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Please select the school you want to view enrollments for"})
			return
		}
	}

	page := pagination.Paginate(c, &paginatedResponse)

	students_, err := enrolmentService.Enrolments(schoolNumber, streamNumber, campusNumber, campusClassNumber, studyPeriodNumber, status, user, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	paginatedResponse.Content = students_

	c.JSON(http.StatusOK, paginatedResponse)
}

func TerminateStudentEnrollment(context *gin.Context) {

	studentNumber := context.Param("student_number")

	// get the enrollment number
	var enrollment models.Enrollment
	if err := config.DB.Where("student_number = ? AND status = ?", studentNumber, constants.StatusActive).First(&enrollment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Student is not enrolled"})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch student enrollment"})
		return
	}

	// update the enrollment status to inactive
	enrollment.Status = constants.StatusTerminated
	if err := config.DB.Save(&enrollment).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to terminate student enrollment"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Student enrollment terminated successfully"})
}

func ReEnrollStudent(context *gin.Context) {
	studentNumber := context.Param("student_number")

	// get the enrollment number
	var enrollment models.Enrollment
	if err := config.DB.Where("student_number = ? AND status = ?", studentNumber, constants.StatusTerminated).First(&enrollment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Student doesnt have a terminated enrollment"})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch student enrollment"})
		return
	}

	// update the enrollment status to active
	enrollment.Status = constants.StatusActive
	if err := config.DB.Save(&enrollment).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to re-enroll student"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Student re-enrolled successfully"})

}
