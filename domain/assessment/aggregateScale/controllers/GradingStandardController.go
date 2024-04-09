package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
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
		Number:            uuid.NewV4().String(),
		Name:              payload.Name,
		Description:       payload.Description,
		Type:              payload.Type,
		CampusClassNumber: payload.ClassNumber,
		UserNumber:        ctx.GetHeader("User-Number"),
	}

	if err := config.DB.Create(&gradingStandard).Error; err != nil {
		ctx.JSON(400, gin.H{"message": "Failed to save grading standard"})
		return
	}

	ctx.JSON(201, gin.H{"message": "Grading standard created successfully", "data": gradingStandard})
}

func FetchGradingStandards(ctx *gin.Context) {
	var gradingStandards []models.GradingStandard
	if err := config.DB.Find(&gradingStandards).Error; err != nil {
		ctx.JSON(400, gin.H{"message": "Failed to get grading standards"})
		return
	}

	ctx.JSON(200, gin.H{"data": gradingStandards})
}

func AttachAssessmentAreas(ctx *gin.Context) {

	var payload requests.AssessmentAreaGradingStandardRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	assessmentAreaGradingStandard := assessmentAreaModel.AssessmentAreaGradingStandard{
		Number:                uuid.NewV4().String(),
		Code:                  payload.Code,
		AssessmentAreaNumber:  payload.AssessmentAreaNumber,
		GradingStandardNumber: payload.GradingStandardNumber,
	}

	if err := config.DB.Create(&assessmentAreaGradingStandard).Error; err != nil {
		ctx.JSON(400, gin.H{"message": "Failed to save assessment area grading standard"})
		return
	}

	ctx.JSON(201, gin.H{"message": "Assessment area grading standard created successfully", "data": assessmentAreaGradingStandard})
}

func FetchAssessmentAreasByGradingStandard(ctx *gin.Context) {
	//gradingStandardNumber := ctx.Param("gradingStandardNumber")

	type response struct {
		AssessmentAreaName  string `json:"assessment_area_name"`
		GradingStandardName string `json:"grading_standard_name"`
	}
	var list []response
	var assessmentAreaGradingStandards []assessmentAreaModel.AssessmentAreaGradingStandard
	if err := config.DB.Find(&assessmentAreaGradingStandards).Error; err != nil {
		ctx.JSON(400, gin.H{"message": "Failed to get assessment area grading standards"})
		return
	}

	for _, val := range assessmentAreaGradingStandards {
		var gradingStandard models.GradingStandard
		var assessmentArea assessmentAreaModel.AssessmentArea
		config.DB.First(&gradingStandard, "number = ?", val.GradingStandardNumber)
		config.DB.First(&assessmentArea, "number = ?", val.AssessmentAreaNumber)
		list = append(list, response{
			AssessmentAreaName:  assessmentArea.Name,
			GradingStandardName: gradingStandard.Name,
		})
	}

	ctx.JSON(200, gin.H{"data": list})

}
