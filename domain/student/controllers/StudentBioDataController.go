package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/request"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/services"
	"github.com/gin-gonic/gin"
)

// UploadStudentBioData handles the file upload for student bio data
func UploadStudentBioData(c *gin.Context) {
	// Parse form data
	err := c.Request.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed to parse form data: " + err.Error(),
		})
		return
	}

	studentNumber := c.PostForm("student_number")
	if studentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "student_number is required",
		})
		return
	}

	// Get the file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed to get file: " + err.Error(),
		})
		return
	}

	// Create request object
	payload := request.StudentBioDataRequest{
		StudentNumber: studentNumber,
		File:          file,
	}

	// Validate request
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed to validate input: " + err.Error(),
		})
		return
	}

	// Upload bio data
	err = services.UploadStudentBioData(c, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Failed to upload student bio data: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Student bio data uploaded successfully",
	})
}

// GetStudentBioData retrieves bio data for a specific student
func GetStudentBioData(c *gin.Context) {
	studentNumber := c.Param("student_number")

	if studentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "student_number is required",
		})
		return
	}

	bioData, err := services.GetStudentBioDataByStudentNumber(studentNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Failed to fetch bio data: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   bioData,
	})
}

// DeleteStudentBioData deletes bio data for a specific student
func DeleteStudentBioData(c *gin.Context) {
	number := c.Param("number")

	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "student_number is required",
		})
		return
	}

	err := services.DeleteStudentBioData(number)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Failed to delete bio data: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Student bio data deleted successfully",
	})
}
