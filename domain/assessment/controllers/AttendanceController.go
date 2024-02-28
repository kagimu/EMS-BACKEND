package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

func CreateAttendance(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AttendanceRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveAttendance := models.Attendance{
		DayOfTheWeek:      payload.DayOfTheWeek,
		CheckIn:           payload.CheckIn,
		CheckOut:          payload.CheckOut,
		ClassNumber:       payload.ClassNumber,
		Present:           payload.Present,
		StudyPeriodNumber: payload.StudyPeriodNumber,
		CampusNumber:      payload.CampusNumber,
		SchoolNumber:      payload.SchoolNumber,
		UserNumber:        user.Number,
	}

	if err := config.DB.Create(&saveAttendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create attendance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": saveAttendance})
}

func FetchStudentAttendance(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	type Payload struct {
		StudentNumber string `json:"student_number"`
		ClassNumber   string `json:"class_number"`
		StudyPeriod   string `json:"study_period"`
		CampusNumber  string `json:"campus_number"`
		SchoolNumber  string `json:"school_number"`
	}

	var payload Payload
	if payload.ClassNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid class number"})
		return
	}
	if payload.StudyPeriod == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid study period"})
		return
	}
	if payload.StudentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid student number"})
		return
	}
	if payload.CampusNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid campus number"})
		return
	}
	if payload.SchoolNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid school number"})
		return
	}

	if err := config.DB.Where("student_number = ?", payload.StudentNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	studentNumber := payload.StudentNumber
	classNumber := payload.ClassNumber
	studyPeriod := payload.StudyPeriod
	campusNumber := payload.CampusNumber
	schoolNumber := payload.SchoolNumber

	var attendance []models.Attendance
	if err := config.DB.Where("student_number = ? AND class_number = ?  AND study_period = ? AND school_number = ? AND campus_number = ?", studentNumber, classNumber, studyPeriod, schoolNumber, campusNumber).Find(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch attendance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": attendance})
}

func FetchAttendanceByClass(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	type Payload struct {
		ClassNumber  string `json:"class_number"`
		StudyPeriod  string `json:"study_period"`
		CampusNumber string `json:"campus_number"`
		SchoolNumber string `json:"school_number"`
	}

	var payload Payload
	if payload.ClassNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid class number"})
		return
	}

	if payload.StudyPeriod == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid study period"})
		return
	}

	classNumber := payload.ClassNumber
	studyPeriod := payload.StudyPeriod
	campusNumber := payload.CampusNumber
	schoolNumber := payload.SchoolNumber

	var attendance []models.Attendance
	if err := config.DB.Where("class_number = ? AND study_period = ? AND campus_number = ? AND school_number = ? ", classNumber, studyPeriod, campusNumber, schoolNumber).Find(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": attendance})
}
