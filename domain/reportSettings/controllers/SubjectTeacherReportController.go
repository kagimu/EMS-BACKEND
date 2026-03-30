package controllers

import (
	"fmt"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/services"
	"github.com/gin-gonic/gin"
)

func CreateSubjectTeacherReport(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	var payload requests.SubjectTeacherReportRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input: %s", err.Error())})
		return
	}
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	err = services.CreateSubjectTeacherReport(payload, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("SubjectTeacherReport created successfully"))
}

func FetchAllSubjectTeacherReports(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	studyPeriodNumber := c.Query("study_period_number")
	classNumber := c.Query("class_number")
	stage := c.Query("stage")
	status := c.Query("status")
	assessmentAreaNumber := c.Query("assessment_area_number")

	reports, err := services.GetSubjectTeacherReports(studyPeriodNumber, stage, status, classNumber, assessmentAreaNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": reports,
	})
}

func UpdateSubjectTeacherReport(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	number := c.Param("number")
	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Provide SubjectTeacherReport number"})
		return
	}

	var payload requests.UpdateSubjectTeacherReportRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input: %s", err.Error())})
		return
	}

	err := services.UpdateSubjectTeacherReport(number, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("SubjectTeacherReport updated successfully"))
}

func RemoveSubjectTeacherReport(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	number := c.Param("number")
	if number == "" {
		c.JSON(http.StatusBadRequest, commons.Feedback("choose SubjectTeacherReport to delete"))
		return
	}

	err := services.DeleteSubjectTeacherReport(number)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("SubjectTeacherReport removed successfully"))
}

func ForwardSubjectTeacherReport(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	var payload requests.ForwardSubjectTeacherReportRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input: %s", err.Error())})
		return
	}

	commenterName := user.FirstName + " " + user.LastName
	err = services.ForwardSubjectTeacherReport(payload, commenterName)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("SubjectTeacherReport forwarded successfully"))
}

func ApproveSubjectTeacherReport(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	number := c.Param("number")
	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Provide SubjectTeacherReport number"})
		return
	}

	err := services.ApproveSubjectTeacherReport(number)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("SubjectTeacherReport approved successfully"))
}

func FetchSubjectTeacherReportComments(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		return
	}

	reportNumber := c.Query("subject_teacher_report_number")
	if reportNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Provide subject_teacher_report_number query param"})
		return
	}

	comments, err := services.GetSubjectTeacherReportComments(reportNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": comments,
	})
}
