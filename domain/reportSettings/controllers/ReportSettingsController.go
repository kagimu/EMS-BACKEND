package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/gin-gonic/gin"

	uuid "github.com/satori/go.uuid"
	"net/http"
)

func CreateReportSettings(c *gin.Context) {
	var payload requests.ReportSettingsRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err := payload.ValidateReportSetting()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	//check for existing report settings with the same study period number, academic year number and campus class number and status active

	var existingReportSettings models.ReportSetting
	if err := config.DB.Where("study_period_number = ? AND academic_year_number = ? AND campus_class_number = ? AND status = ?", payload.StudyPeriodNumber, payload.AcademicYearNumber, payload.CampusClassNumber, constants.STATUS_ACTIVE).First(&existingReportSettings).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Report settings for this class  already exists"})
		return

	}

	// Save report settings
	reportSettings := models.ReportSetting{
		Number:             uuid.NewV4().String(),
		StudyPeriodNumber:  payload.StudyPeriodNumber,
		AcademicYearNumber: payload.AcademicYearNumber,
		CampusClassNumber:  payload.CampusClassNumber,
		Columns:            payload.Columns,
		Status:             constants.STATUS_ACTIVE,
		UserNumber:         c.Request.Header["User-Number"][0],
	}

	if err := config.DB.Create(&reportSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create report settings" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report settings created successfully"})

}

func FetchReportSettings(c *gin.Context) {
	var reportSettings []models.ReportSetting
	if err := config.DB.Preload("AcademicYear").Preload("CampusClass").Preload("StudyPeriod").Find(&reportSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report settings"})
		return
	}

	var rs []responses.ReportSettingsResponse

	for _, reportSetting := range reportSettings {
		response := responses.ReportSettingsResponse{
			Name:    reportSetting.AcademicYear.AcademicYearName + " - " + reportSetting.CampusClass.ClassName + " - " + reportSetting.StudyPeriod.StudyPeriodName,
			Number:  reportSetting.Number,
			Status:  reportSetting.Status,
			Columns: reportSetting.Columns,
		}

		rs = append(rs, response)
	}

	c.JSON(http.StatusOK, gin.H{"data": rs})

}

func CreatReportColumn(ctx *gin.Context) {
	var payload requests.ReportColumnRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err := payload.ValidateReportColumn()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Save report column
	reportColumn := models.ReportColumn{
		Number:              uuid.NewV4().String(),
		ReportSettingNumber: payload.ReportSettingNumber,
		Name:                payload.Name,
		PercentageShare:     payload.PercentageShare,
		UserNumber:          ctx.Request.Header["User-Number"][0],
	}

	if err := config.DB.Create(&reportColumn).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create report column"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Report column created successfully"})

}

func FetchReportColumns(ctx *gin.Context) {
	var reportColumns []models.ReportColumn
	if err := config.DB.Find(&reportColumns).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report columns"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": reportColumns})

}

func CreateReportColumnAssessment(ctx *gin.Context) {
	var payload requests.ReportColumnAssessmentRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err := payload.ValidateReportColumnAssessment()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Save report column assessment
	reportColumnAssessment := models.ReportColumnAssessment{
		Number:                 uuid.NewV4().String(),
		ReportColumnNumber:     payload.ReportColumnNumber,
		AssessmentAreaNumber:   payload.AssessmentAreaNumber,
		AssessmentMetricNumber: payload.AssessmentMetricNumber,
		UserNumber:             ctx.Request.Header["User-Number"][0],
	}

	if err := config.DB.Create(&reportColumnAssessment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create report column assessment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Report column assessment created successfully"})
}

func FetchReportColumnAssessments(ctx *gin.Context) {
	var reportColumnAssessments []models.ReportColumnAssessment
	if err := config.DB.Find(&reportColumnAssessments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report column assessments"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": reportColumnAssessments})

}

func CreateReportSummary(ctx *gin.Context) {
	var payload requests.ReportSummaryRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err := payload.ValidateReportSummary()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Save report summary
	reportSummary := models.ReportSummary{
		Number:              uuid.NewV4().String(),
		ReportSettingNumber: payload.ReportSettingNumber,
		StudentNumber:       payload.StudentNumber,
		CampusClassNumber:   payload.CampusClassNumber,
		StudyPeriodNumber:   payload.StudyPeriodNumber,
		Grade:               payload.Grade,
		Aggregate:           payload.Aggregate,
		Comment:             payload.Comment,
		UserNumber:          ctx.Request.Header["User-Number"][0],
	}

	if err := config.DB.Create(&reportSummary).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create report summary"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Report summary created successfully"})
}

func FetchReportSummaries(ctx *gin.Context) {
	var reportSummaries []models.ReportSummary
	if err := config.DB.Find(&reportSummaries).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report summaries"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": reportSummaries})

}

func CreateReportScore(ctx *gin.Context) {
	var payload requests.ReportScoreRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err := payload.ValidateReportScore()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Save report score
	reportScore := models.ReportScore{
		ReportSummaryNumber:  payload.ReportSummaryNumber,
		AssessmentAreaNumber: payload.AssessmentAreaNumber,
		ReportColumnNumber:   payload.ReportColumnNumber,
		Score:                payload.Score,
		Aggregate:            payload.Aggregate,
		Grade:                payload.Grade,
		UserNumber:           ctx.Request.Header["User-Number"][0],
	}

	if err := config.DB.Create(&reportScore).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create report score"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Report score created successfully"})
}

func FetchReportScores(ctx *gin.Context) {
	var reportScores []models.ReportScore
	if err := config.DB.Find(&reportScores).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report scores"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": reportScores})

}

func CreateReportScoreSummary(ctx *gin.Context) {
	var payload requests.ReportScoreSummaryRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err := payload.ValidateReportScoreSummary()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Save report score summary
	reportScoreSummary := models.ReportScoreSummary{
		ReportSummaryNumber:  payload.ReportSummaryNumber,
		AssessmentAreaNumber: payload.AssessmentAreaNumber,
		Aggregate:            payload.Aggregate,
		Grade:                payload.Grade,
		Summation:            payload.Summation,
		UserNumber:           ctx.Request.Header["User-Number"][0],
	}

	if err := config.DB.Create(&reportScoreSummary).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create report score summary"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Report score summary created successfully"})
}

func FetchReportScoreSummaries(ctx *gin.Context) {
	var reportScoreSummaries []models.ReportScoreSummary
	if err := config.DB.Find(&reportScoreSummaries).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report score summaries"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"report_score_summaries": reportScoreSummaries})

}
