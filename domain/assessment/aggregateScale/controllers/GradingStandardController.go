package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateGradingStandard(ctx *gin.Context) {
	var payload requests.GradingStandardRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	// Save GradingStandard
	gradingStandard := models.GradingStandard{
		Number:      uuid.NewV4().String(),
		Name:        payload.Name,
		Description: payload.Description,
		Type:        payload.Type,
		ClassNumber: payload.ClassNumber,
		UserNumber:  ctx.GetHeader("User-Number"),
	}

	if err := config.DB.Create(&gradingStandard).Error; err != nil {
		ctx.JSON(400, gin.H{"message": "Failed to save grading standard"})
		return
	}

	ctx.JSON(201, gin.H{"message": "Grading standard created successfully", "data": gradingStandard})
}

func GetGradingStandards(ctx *gin.Context) {
	var gradingStandards []models.GradingStandard
	if err := config.DB.Find(&gradingStandards).Error; err != nil {
		ctx.JSON(400, gin.H{"message": "Failed to get grading standards"})
		return
	}

	ctx.JSON(200, gin.H{"data": gradingStandards})
}
