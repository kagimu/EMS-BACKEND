package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
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

	if err := services.CreateGrade(ctx, request); err != nil {
		ctx.JSON(500, gin.H{"message": "Failed to create grade" + err.Error()})
		return
	}

	ctx.JSON(201, commons.Feedback("Grade created successfully"))

}

func FetchGrades(ctx *gin.Context) {
	grades, err := services.FetchGrades(ctx)

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

	if err := services.UpdateGrade(ctx, request); err != nil {
		ctx.JSON(500, gin.H{"message": "Failed to update grade" + err.Error()})
		return
	}

	ctx.JSON(200, commons.Feedback("Grade updated successfully"))
}
