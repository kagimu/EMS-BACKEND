package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/request"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateStudent(c *gin.Context) {
	var payload request.StudentRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Internal server error!",
		})
		return
	}

	payload.Prepare()
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to validate input",
		})
		return
	}

	// change this
	userNumber := "b7426afb-d114-479b-b83b-4ec59a2bbe64"
	schoolNumber := c.Param("school_number")
	classNumber := c.Param("class_number")
	sourceOfSchoolInfoNumber := c.Param("source_of_school_info_number")

	apnGenerator := models.NewAPNGenerator()

	studentAPN, err := apnGenerator.GenerateAPN()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to generate APN",
			"error":   err,
		})

		return
	}

	apn := studentAPN.APN
	alternativeAPN := fmt.Sprintf("%d%05d", studentAPN.Year, studentAPN.Serial)

	var user authModels.User
	if err := config.DB.
		Where("number = ?", userNumber).
		First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	tx := config.DB.Begin()
	saveStudent := models.Student{
		Number:                   uuid.NewV4().String(),
		APN:                      apn,
		APNAlternative:           alternativeAPN,
		Surname:                  payload.Surname,
		OtherName:                payload.OtherName,
		Age:                      payload.Age,
		Gender:                   payload.Gender,
		Reason:                   payload.Reason,
		SecondName:               payload.SecondName,
		UserNumber:               userNumber,
		SchoolNumber:             schoolNumber,
		ClassNumber:              classNumber,
		SourceOfSchoolInfoNumber: sourceOfSchoolInfoNumber,
	}

	if err := tx.Create(&saveStudent).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "failed to save source of school info",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to save school of school info"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "source of school information created successfully",
		"data":    saveStudent,
	})
}

func UpdateStudent(c *gin.Context) {
	student := models.Student{}
	studentAPN := c.Param("id")

	err := config.DB.First(&student, studentAPN).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no student found",
		})
		return
	}

	if err := c.BindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "invalid request",
		})
		return
	}

	if err := config.DB.Save(&student).Error; err != nil {
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

func FetchStudent(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": student})
}

func FetchAllStudents(c *gin.Context) {
	var students []models.Student

	err := config.DB.Find(&students).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to fetch students",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": students})
}
