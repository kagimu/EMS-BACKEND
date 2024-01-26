package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/request"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreateStudent(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found in header",
		})
		return
	}
	var payload request.StudentRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	payload.Prepare()
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("failed to validate input %s:", err.Error()),
		})
		return
	}

	applicationNumber := payload.ApplicationNumber
	var application applicationModel.Application
	if err := config.DB.Where("number = ?", applicationNumber).First(&application).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Application not found"})
		return
	}

	studentAPN, err := CreateAPN()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create APN"})
		return
	}

	studentNumber := uuid.NewV4().String()
	userNumber := user.Number
	savedStudent := models.Student{
		Number:                    studentNumber,
		APN:                       studentAPN,
		Surname:                   payload.Surname,
		OtherName:                 payload.OtherName,
		Gender:                    payload.Gender,
		DateOfBirth:               payload.DateOfBirth,
		ReasonForChoosingSchool:   payload.ReasonForChoosingSchool,
		PreferredTypeOfEnrollment: payload.PreferredTypeOfEnrollment,
		SecondName:                payload.SecondName,
		PreferedCampusNumber:      payload.PreferredCampusNumber,
		UserNumber:                userNumber,
		SourceOfSchoolInfo:        payload.SourceOfSchoolInfo,
		ApplicationNumber:         applicationNumber,
		ClassNumber:               payload.ClassNumber,
	}

	tx := config.DB.Begin()

	if err := tx.Create(&savedStudent).Error; err != nil {
		tx.Rollback()

		c.JSON(http.StatusBadRequest, gin.H{
			"details": err.Error(),
			"message": "failed to save student",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save student"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "student created successfully",
		"data":    savedStudent,
	})
}

func EditStudent(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	var payload request.StudentUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to process input",
		})
		return
	}

	payload.Prepare()

	var student models.Student
	studentAPN := c.Param("apn")
	if err := config.DB.First(&student, "apn = ?", studentAPN).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No student found with the specified APN",
		})
		return
	}

	savedStudent := models.Student{
		Surname:                   payload.Surname,
		OtherName:                 payload.OtherName,
		SecondName:                payload.SecondName,
		APNStatus:                 payload.APNStatus,
		DateOfBirth:               payload.DateOfBirth,
		ReasonForChoosingSchool:   payload.ReasonForChoosingSchool,
		PreferredTypeOfEnrollment: payload.PreferredTypeOfEnrollment,
		ClassNumber:               payload.ClassNumber,
	}

	if err := config.DB.Model(&student).Where("apn = ?", studentAPN).Updates(savedStudent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to update student",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": student})
}

func DeleteStudent(c *gin.Context) {
	studentAPN := c.Param("id")
	student := models.Student{}

	err := config.DB.First(&student, studentAPN).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no student found",
		})
		return
	}

	if err := config.DB.Delete(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to delete student",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "student deleted successfully"})
}

func FetchStudentByAPN(c *gin.Context) {
	if user, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"debug":   user,
			"message": "User not found"})
		return
	}

	var student models.Student
	studentAPN := c.Param("apn")

	if studentAPN == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Student API is missing",
		})
		return
	}

	if err := config.DB.First(&student, "apn = ?", studentAPN).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "no student found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": student})
}

func FetchAllStudents(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var students []models.Student

	if err := config.DB.Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to fetch students",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": students})
}
