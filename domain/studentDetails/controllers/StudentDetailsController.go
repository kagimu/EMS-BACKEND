package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentDetails/api/requests"
	studentDetails "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentDetails/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateStudentDetails(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "User is not found"})
		return
	}

	var payload requests.StudentDetailsRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("failed to validate input %s:", err.Error())})
		return
	}

	studentNumber := c.Param("student_number")
	var student models.Student
	if err := config.DB.Where("number = ?", studentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"debug":   err.Error(),
			"message": "Student not found"})
		return
	}

	tx := config.DB.Begin()
	userNumber := user.Number

	saveStudentDetails := studentDetails.StudentDetail{
		Number:             uuid.NewV4().String(),
		Surname:            payload.Surname,
		SecondName:         payload.SecondName,
		OtherName:          payload.OtherName,
		DateOfBirth:        payload.DateOfBirth,
		BirthCertificateNo: payload.BirthCertificateNo,
		CountryOfBirth:     payload.CountryOfBirth,
		HomeDistrict:       payload.HomeDistrict,
		HomeSubCounty:      payload.HomeSubCounty,
		HomeVillage:        payload.HomeVillage,
		HomePhone1:         payload.HomePhone1,
		HomePhone2:         payload.HomePhone2,
		HomeEmail1:         payload.HomeEmail1,
		HomeEmail2:         payload.HomeEmail2,
		Religion:           payload.Religion,
		HomeCountry:        payload.HomeCountry,
		HomeCounty:         payload.HomeCounty,
		HomeParish:         payload.HomeParish,
		NearestTown:        payload.NearestTown,
		StudentNumber:      studentNumber,
		UserNumber:         userNumber,
	}

	if err := tx.Create(&saveStudentDetails).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create  Student details"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to add student details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Student details Added Successfully!", "data": saveStudentDetails})
}

func EditStudentDetails(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	var payload requests.StudentDetailsRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	var studentDetailsModel studentDetails.StudentDetail
	studentNumber := c.Param("student_number")

	if err := config.DB.Where("number = ?", studentNumber).First(&studentDetailsModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	updateStudentDetails := studentDetails.StudentDetail{
		Surname:            payload.Surname,
		SecondName:         payload.SecondName,
		OtherName:          payload.OtherName,
		DateOfBirth:        payload.DateOfBirth,
		BirthCertificateNo: payload.BirthCertificateNo,
		CountryOfBirth:     payload.CountryOfBirth,
		HomeDistrict:       payload.HomeDistrict,
		HomeSubCounty:      payload.HomeSubCounty,
		HomeVillage:        payload.HomeVillage,
		HomePhone1:         payload.HomePhone1,
		HomePhone2:         payload.HomePhone2,
		HomeEmail1:         payload.HomeEmail1,
		HomeEmail2:         payload.HomeEmail2,
		Religion:           payload.Religion,
		HomeCountry:        payload.HomeCountry,
		HomeCounty:         payload.HomeCounty,
		HomeParish:         payload.HomeParish,
		NearestTown:        payload.NearestTown,
	}

	if err := config.DB.Model(&studentDetailsModel).Where("number = ?", studentNumber).Updates(updateStudentDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to update student details",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": studentDetailsModel})
}
