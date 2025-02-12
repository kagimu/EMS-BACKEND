package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/services"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
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

	checkinTime, err := commons.ConvertStringToDate(payload.CheckInTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to convert check in time"})
		return
	}

	checkOutTime, err := commons.ConvertStringToDate(payload.CheckOutTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to convert check out time"})
		return
	}
	date, err := commons.ConvertStringToDate(payload.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to convert date"})
		return
	}

	saveAttendance := models.Attendance{
		Present:           payload.Present,
		CheckInTime:       checkinTime,
		CheckOutTime:      checkOutTime,
		Date:              date,
		SchoolNumber:      payload.SchoolNumber,
		UserNumber:        user.Number,
		CampusClassNumber: payload.CampusClassNumber,
		CampusNumber:      payload.CampusNumber,
		AddedBy:           user.Number,
	}

	if err := services.RecordAttendance(&saveAttendance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create attendance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": saveAttendance})
}

func FetchAttendance(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	//studentNumber := c.Query("student_number")
	//classNumber := c.Query("campus_class_number")
	//studyPeriodNumber := c.Query("study_period")
	//campusNumber := c.Query("campus_number")
	//schoolNumber := c.Query("school_number")

	attendance, err := services.FetchAttendance(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch attendance records"})
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
