package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
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

	date, err := commons.ConvertStringToDate(payload.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to convert date"})
		return
	}

	saveAttendance := models.Attendance{
		Number:            uuid.NewV4().String(),
		Present:           payload.Present,
		CheckInTime:       payload.CheckInTime,
		CheckOutTime:      payload.CheckOutTime,
		Date:              date,
		SchoolNumber:      payload.SchoolNumber,
		UserNumber:        payload.UserNumber,
		CampusClassNumber: payload.CampusClassNumber,
		CampusNumber:      payload.CampusNumber,
		AddedBy:           user.Number,
		StudentNumber:     payload.StudentNumber,
		Comment:           payload.AbsenceReason,
		StreamNumber:      payload.StreamNumber,
		StudyPeriodNumber: payload.StudyPeriodNumber,
	}

	if err := services.RecordAttendance(&saveAttendance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create attendance" + err.Error()})
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

func CheckOut(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.CheckoutRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.ValidateCheckout(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	repo := repositories.NewAttendanceRepository(config.DB)
	attendance, err := repo.GetAttendanceByNumber(payload.Number)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Failed to fetch attendance"})
		return
	}

	attendance.CheckOutTime = payload.CheckoutTime
	attendance.CheckedOutBy = user.Number
	attendance.CheckoutComment = payload.Comment

	err = services.RecordCheckout(&attendance)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Failed to record  checkout" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checkout saved successfully"})
}
