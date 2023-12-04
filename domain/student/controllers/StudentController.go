package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
)

func CreateStudent(c *gin.Context) {
	var student models.Student
	var err error

	err = c.BindJSON(&student)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	student.Prepare()

	err = student.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to validate input",
		})
		return
	}

	apnGenerator := models.NewAPNGenerator()

	studentAPN, err := apnGenerator.GenerateAPN()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to generate APN",
			"error":   err,
		})
		return
	}

	student.APN = studentAPN.APN

	err = config.DB.Create(&student).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to save student info",
			"error":   err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "student created successfully",
		"student": student,
	})
}

func UpdateStudent(c *gin.Context) {
	student := models.Student{}
	studentAPN := c.Param("apn")

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
	studentAPN := c.Param("apn")
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
	studentAPN := c.Param("apn")
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
