package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
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
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	studentNumber := c.Query("student_number")
	classNumber := c.Query("class_number")
	studyPeriodNumber := c.Query("study_period")
	campusNumber := c.Query("campus_number")
	schoolNumber := c.Query("school_number")

	var student studentModel.Student
	if err := config.DB.Where("number = ?", studentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	var class classModel.Class
	if err := config.DB.Where("number = ?", classNumber).First(&class).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Class not found"})
		return
	}

	var studyPeriod studyPeriodModel.StudyPeriod
	if err := config.DB.Where("number = ?", studyPeriodNumber).First(&studyPeriod).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Study period not found"})
		return
	}

	var campus campusModel.Campus
	if err := config.DB.Where("number = ?", studyPeriodNumber).First(&campus).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Campus not found"})
		return
	}

	var school schoolModel.School
	if err := config.DB.Where("number = ?", schoolNumber).First(&school).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "School not found"})
		return
	}

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

	classNumber := c.Query("class_number")
	studyPeriodNumber := c.Query("study_period")
	campusNumber := c.Query("campus_number")
	schoolNumber := c.Query("school_number")

	var class classModel.Class
	if err := config.DB.Where("number = ?", classNumber).First(&class).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Class not found"})
		return
	}

	var studyPeriod studyPeriodModel.StudyPeriod
	if err := config.DB.Where("number = ?", studyPeriodNumber).First(&studyPeriod).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Study period not found"})
		return
	}

	var campus campusModel.Campus
	if err := config.DB.Where("number = ?", studyPeriodNumber).First(&campus).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Campus not found"})
		return
	}

	var school schoolModel.School
	if err := config.DB.Where("number = ?", schoolNumber).First(&school).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "School not found"})
		return
	}

	var attendance []models.Attendance
	if err := config.DB.Where("class_number = ? AND study_period = ? AND campus_number = ? AND school_number = ? ", classNumber, studyPeriod, campusNumber, schoolNumber).Find(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch attendance records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": attendance})
}
