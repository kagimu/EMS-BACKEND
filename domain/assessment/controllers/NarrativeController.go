package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var narrativeService = services.NewNarrativeService()

func SaveNarrative(ctx *gin.Context) {

	user, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var request requests.NarrativeRequest
	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = narrativeService.CreateNarrative(request, user.Number)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Narrative saved"})
}

func GetNarrative(context *gin.Context) {

	_startDate := context.Query("start_date")
	_endDate := context.Query("end_date")

	var startDate time.Time
	var endDate time.Time
	var err error

	if _startDate != "" {
		startDate, err = time.Parse(time.DateOnly, _startDate)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Error reading start date: " + err.Error()})
			return
		}
	}

	if _endDate != "" {
		endDate, err = time.Parse(time.DateOnly, _endDate)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Error reading end date: " + err.Error()})
			return
		}
	}

	narrative := requests.NarrativeRequest{
		Description:          context.Query("description"),
		StudentNumber:        context.Query("student_number"),
		AssessmentAreaNumber: context.Query("assessment_area_number"),
		StudyPeriodNumber:    context.Query("study_period_number"),
		CampusClassNumber:    context.Query("campus_class_number"),
	}

	data, err := narrativeService.AllNarrative(narrative, startDate, endDate)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": data})
}

func SaveNarrativeSchedule(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var request requests.ScheduledNarrativeRequest
	err = context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = narrativeService.SaveSchedule(request, user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Scheduled Narrative saved"})
}

func NarrativeSchedule(context *gin.Context) {
	studyPeriodNumber := context.Query("study_period_number")
	classNumber := context.Query("class_number")
	campusNumber := context.Query("campus_number")
	curriculumNumber := context.Query("curriculum_number")
	assessmentAreaNumber := context.Query("assessment_area_number")

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "choose class"})
		return
	}

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "choose campus"})
		return
	}

	if curriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "choose curriculum"})
		return
	}

	if studyPeriodNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "study_period_number is required"})
		return
	}

	data, err := narrativeService.AllSchedules(studyPeriodNumber, classNumber, campusNumber, curriculumNumber, assessmentAreaNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": data})
}

func GetDetailedNarrative(context *gin.Context) {
	narrativeNumber := context.Query("narrative_number")
	if narrativeNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "narrative number is required"})
		return
	}

	data, err := narrativeService.NarrativeDetails(narrativeNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": data})
}

func DeleteNarrative(context *gin.Context) {
	narrativeNumber := context.Query("narrative_number")
	if narrativeNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "narrative number is required"})
		return
	}

	err := narrativeService.DeleteNarrative(narrativeNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Removed successfully"})
}

func DeleteNarrativeSchedule(context *gin.Context) {

}

func UpdateNarrative(ctx *gin.Context) {
	narrativeNumber := ctx.Param("narrative_number")

	if narrativeNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "narrative number is required"})
		return
	}

	user, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var request requests.NarrativeRequest
	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = narrativeService.UpdateNarrativeEntry(request, narrativeNumber, user.Number)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Entry updated successfully"})
}
