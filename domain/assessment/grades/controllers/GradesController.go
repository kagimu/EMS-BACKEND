package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateGrade(ctx *gin.Context) {

	var request requests.GradeDataRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"message": "Failed to validate input"})
		return
	}

	if err := request.Validate(); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	userNumber := ctx.GetHeader("User-Number")

	if err := services.CreateGrade(userNumber, request); err != nil {
		ctx.JSON(400, gin.H{"message": "Failed to create grade" + err.Error()})
		return
	}

	ctx.JSON(201, commons.Feedback("Name of aggregate created successfully"))

}

func FetchGrades(ctx *gin.Context) {

	schoolNumber := ctx.Query("school_number")

	grades, err := services.FetchGrades(schoolNumber, ctx)

	if err != nil {
		ctx.JSON(500, gin.H{"message": "Failed to fetch grades"})
		return
	}

	ctx.JSON(200, grades)

}

func UpdateGrade(ctx *gin.Context) {

	var request requests.EditGradeDataRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"message": "Failed to validate input"})
		return
	}

	if err := request.Validate(); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	userNumber := ctx.GetHeader("User-Number")

	if err := services.UpdateGrade(userNumber, request); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, commons.Feedback("Grade updated successfully"))
}

func RemoveGrade(context *gin.Context) {
	gradeNumber := context.Param("grade_number")
	if gradeNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Provide grade number to delete"})
		return
	}

	err := services.RemoveGrade(gradeNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Grade removed successfully"})
}
