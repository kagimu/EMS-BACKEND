package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	assessmentAreaServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateGradingStandard(ctx *gin.Context) {
	user, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get user from header"})
		return
	}

	var payload requests.GradingStandardRequest
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.SaveGradingStandard(payload, user.Number)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Grading standard created successfully"})
}

func UpdateGradingStandard(ctx *gin.Context) {
	user, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get user from header"})
		return
	}

	var payload requests.UpdateGradingStandardRequest
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.UpdateGradingStandard(payload, user.Number)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Grading standard updated successfully"})
}

func FetchGradingStandards(ctx *gin.Context) {
	schoolNumber := ctx.Query("school_number")
	classNumber := ctx.Query("class_number")
	scaleNumber := ctx.Query("scale_number")

	response, err := services.FetchGradingStandards(schoolNumber, classNumber, scaleNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": response})
}

func AttachAssessmentAreas(ctx *gin.Context) {

	var payload requests.AssessmentAreaGradingStandardRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := assessmentAreaServices.SaveAssessmentAreaGradingStandard(payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "Assessment area grading standard created successfully"})
}

func FetchAssessmentAreasByGradingStandard(ctx *gin.Context) {

	gradingStandardNumber := ctx.Query("grading_standard_number")
	assessmentAreaNumber := ctx.Query("assessment_area_number")
	schoolNumber := ctx.Query("school_number")
	if schoolNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "choose a school"})
		return
	}

	response, err := assessmentAreaServices.FetchAssessmentAreasByGradingStandard(schoolNumber, assessmentAreaNumber, gradingStandardNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})

}

func RemoveGradingStandard(context *gin.Context) {
	gradingStandardNumber := context.Param("grading_standard_number")
	if gradingStandardNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Provide grading standard number"})
		return
	}

	err := services.RemoveGradingStandard(gradingStandardNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, commons.Feedback("Removed grading standard"))
}

func DetailsGradingStandard(context *gin.Context) {

	gradingStandardNumber := context.Query("grading_standard_number")
	if gradingStandardNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Provide grading standard number"})
		return
	}

	detailedResponse, err := services.GetGradingStandardDetails(gradingStandardNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": detailedResponse})

}

func UpdateStatusGradingStandard(ctx *gin.Context) {

	user, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get user from header"})
		return
	}

	var payload requests.UpdateGradingStatusStandardRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.UpdateGradingStandardStatus(payload, user.Number)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "Status updated successfully"})
}

func MakeGradingStandardDefault(ctx *gin.Context) {

	var payload requests.MakeGradingStandardDefault
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := services.MakeGradingStandardDefault(payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "Default grading standard set successfully"})
}

func ClassGradingStandards(context *gin.Context) {
	schoolNumber := context.Query("school_number")
	gradingStandardNumber := context.Query("grading_standard_number")
	classNumber := context.Query("class_number")
	scaleNumber := context.Query("scale_number")

	response, err := services.ClassGradingStandards(schoolNumber, gradingStandardNumber, classNumber, scaleNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(200, gin.H{"data": response})
}
