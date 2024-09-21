package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

var billableService = services.NewBillableService()

func CreateBillable(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.BillableRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input %s:", err.Error())})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err = billableService.NewBillable(payload, user.Number); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Billable created successfully"})
}

func AllBillable(context *gin.Context) {
	schoolNumber := context.Query("school_number")
	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose school"})
		return
	}

	status := context.Query("status")
	campusNumber := context.Query("campus_number")
	billType := context.Query("type")
	classNumber := context.Query("class_number")
	studyPeriodNumber := context.Query("study_period_number")
	academicYearNumber := context.Query("academic_year_number")

	all, err := billableService.FindAllBillable(schoolNumber, campusNumber, billType, classNumber, studyPeriodNumber, academicYearNumber, status)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "All billable", "data": all})
}

func BillableDetails(context *gin.Context) {
	billableNumber := context.Query("billable_number")
	if billableNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose billable item"})
		return
	}

	billable, err := billableService.BillableDetails(billableNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": billable})
}

func GetBill(context *gin.Context) {
	studentNumber := context.Query("student_number")
	studyPeriodNumber := context.Query("study_period_number")
	classNumber := context.Query("class_number")
	modeOfStudyNumber := context.Query("mode_of_study_number")
	//curriculumNumber := context.Query("curriculum_number")
	campusCurriculumNumber := context.Query("campus_curriculum_number")
	campusNumber := context.Query("campus_number")
	if studentNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose student"})
		return
	}

	if studyPeriodNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose study period"})
		return
	}

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose class"})
		return
	}

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose campus"})
		return
	}

	if modeOfStudyNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose mode of study"})
		return
	}

	billable, err := billableService.GetBill(studentNumber, campusNumber, classNumber, modeOfStudyNumber, studyPeriodNumber, campusCurriculumNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": billable})
}

func DeleteBillable(context *gin.Context) {

	if _, err := commons.GetUserFromHeader(context); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	billableNumber := context.Param("billable_number")

	if billableNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose bill to delete"})
		return
	}

	err := billableService.DeleteBillable(billableNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, commons.Feedback("Billable removed successfully"))
}

func UpdateBillable(context *gin.Context) {
	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.UpdateBillableRequest
	if err := context.ShouldBindJSON(&payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Failed to validate input %s:", err.Error())})
		return
	}

	if err := payload.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err = billableService.UpdateBillable(payload, user.Number); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Billable updated successfully"})
}
