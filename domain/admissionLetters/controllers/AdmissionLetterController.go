package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admissionLetters/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admissionLetters/services"
	"github.com/gin-gonic/gin"
)

var admissionLetterService = services.NewAdmissionLetterService()

func CreateAdmissionLetter(c *gin.Context) {
	var req requests.AdmissionLetterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	data, err := admissionLetterService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Admission letter created successfully", "data": data})
}

func UpdateAdmissionLetter(c *gin.Context) {
	var req requests.UpdateAdmissionLetterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	req.Number = c.Param("number")
	if req.Number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "admission letter number is required"})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	data, err := admissionLetterService.Update(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Admission letter updated successfully", "data": data})
}

func GetAdmissionLetter(c *gin.Context) {
	number := c.Param("number")
	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "admission letter number is required"})
		return
	}
	data, err := admissionLetterService.GetByNumber(number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": data})
}

func GetAdmissionLettersByStudent(c *gin.Context) {
	studentNumber := c.Query("student_number")
	if studentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "student_number is required"})
		return
	}
	data, err := admissionLetterService.GetByStudentNumber(studentNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": data})
}

func GetAdmissionLettersBySchool(c *gin.Context) {
	schoolNumber := c.Query("school_number")
	if schoolNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "school_number is required"})
		return
	}
	data, err := admissionLetterService.GetBySchoolNumber(schoolNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": data})
}

func DeleteAdmissionLetter(c *gin.Context) {
	number := c.Param("number")
	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "admission letter number is required"})
		return
	}
	if err := admissionLetterService.Delete(number); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Admission letter deleted successfully"})
}
