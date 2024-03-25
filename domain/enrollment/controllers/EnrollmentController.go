package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	admissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
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

	// validate if student exists on the same campus and school and campus_class_number and study period and   is active
	var existingEnrollment models.Enrollment

	if err := config.DB.Where("student_number = ? AND campus_number = ? AND school_number = ? AND campus_class_number = ? AND study_period_number = ? AND status = ?", payload.StudentNumber, payload.CampusNumber, payload.SchoolNumber, payload.CampusClassNumber, payload.StudyPeriodNumber, "active").First(&existingEnrollment).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student is already enrolled in this class"})
		return
	}

	saveEnrollment := models.Enrollment{
		Number:            uuid.NewV4().String(),
		StudentNumber:     payload.StudentNumber,
		CampusNumber:      payload.CampusNumber,
		SchoolNumber:      payload.SchoolNumber,
		EnrollmentDate:    payload.EnrollmentDate,
		StudyPeriodNumber: payload.StudyPeriodNumber,
		CampusClassNumber: payload.CampusClassNumber,

		Status:     payload.Status,
		UserNumber: user.Number,
	}

	tx := config.DB.Begin()

	if err := tx.Create(&saveEnrollment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save enrollment"})

		return
	}

	// update admission status to enrolled where student number is equal to payload.StudentNumber

	var admission admissionModel.Admission

	if err := tx.Where("student_number = ?", payload.StudentNumber).First(&admission).Error; err != nil {
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"message": "Enrollment created successfully"})
		return
	}

	admission.Status = constants.STATUS_ENROLLED

	if err := tx.Save(&admission).Error; err != nil {

	}

	tx.Commit()

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
