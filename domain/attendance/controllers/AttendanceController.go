package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/services"
	"net/http"
	"time"

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
	var checkinTime, checkOutTime, date time.Time

	if payload.CheckInTime == "" {
		checkinTime = time.Now()
	} else {
		checkinTime, err = commons.ConvertStringToDate(payload.CheckInTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to convert check in time"})
			return
		}
	}
	if payload.CheckOutTime == "" {
		checkOutTime = time.Now()
	} else {
		checkOutTime, err = commons.ConvertStringToDate(payload.CheckOutTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to convert check out time"})
			return
		}
	}

	date, err = commons.ConvertStringToDate(payload.Date)
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
