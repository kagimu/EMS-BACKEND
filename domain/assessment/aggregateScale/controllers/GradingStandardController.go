package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
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
	campusClassNumber := ctx.Query("campus_class_number")
	campusNumber := ctx.Query("campus_number")

	var gss = make([]responses.GradingStandardResponse, 0)
	if campusClassNumber == "" {
		ctx.JSON(http.StatusOK, gin.H{"data": gss})
		return
	}
	var gradingStandards []models.GradingStandard
	if err := config.DB.Scopes(filters.ByCampusClassNumber(campusClassNumber),
		filters.ByCampusNumber(campusNumber)).Preload("Campus").Preload("CampusClass").Find(&gradingStandards).Error; err != nil {
		ctx.JSON(400, gin.H{"message": "Failed to get grading standards"})
		return
	}

	var gs []responses.GradingStandardResponse

	for _, val := range gradingStandards {
		gs = append(gs, responses.GradingStandardResponse{
			Code:        val.Code,
			Number:      val.Number,
			Name:        val.Name,
			Description: val.Description,
			Type:        val.Type,
			Campus:      val.Campus.Name,
			Class:       val.CampusClass.ClassName,
		})
	}

	ctx.JSON(200, gin.H{"data": gs})
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
		Code                int64  `json:"code"`
		AssessmentAreaName  string `json:"assessment_area_name"`
		GradingStandardName string `json:"grading_standard_name"`
	}
	var list []response
	var assessmentAreaGradingStandards []assessmentAreaModel.AssessmentAreaGradingStandard
	if err := config.DB.Scopes(filters.ByGradingStandardNumber(ctx.Query("grading_standard_number"))).Find(&assessmentAreaGradingStandards).Error; err != nil {
		ctx.JSON(400, gin.H{"message": "Failed to get assessment area grading standards"})
		return
	}

	for _, val := range assessmentAreaGradingStandards {
		var gradingStandard models.GradingStandard
		var assessmentArea assessmentAreaModel.AssessmentArea
		config.DB.First(&gradingStandard, "number = ?", val.GradingStandardNumber)
		config.DB.First(&assessmentArea, "number = ?", val.AssessmentAreaNumber)
		list = append(list, response{
			Code:                val.Code,
			AssessmentAreaName:  assessmentArea.Name,
			GradingStandardName: gradingStandard.Name,
		})
	}

	ctx.JSON(200, gin.H{"data": list})

}
