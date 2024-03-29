package controllers

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assesmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	enrollmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"

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
		Code:               payload.Code,
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
			Code:    reportSetting.Code,
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
		Code:                payload.Code,
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

	var rc []responses.ReportColumnResponse

	for _, reportColumn := range reportColumns {
		response := responses.ReportColumnResponse{
			Name:            reportColumn.Name,
			Number:          reportColumn.Number,
			PercentageShare: reportColumn.PercentageShare,
			Code:            reportColumn.Code,
		}
		rc = append(rc, response)

	}

	ctx.JSON(http.StatusOK, gin.H{"data": rc})

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
		Code:                   payload.Code,
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

	var rca []responses.ReportColumnAssessmentResponse

	for _, reportColumnAssessment := range reportColumnAssessments {
		var assessmentArea assessmentAreaModel.AssessmentArea
		if err := config.DB.Where("number = ?", reportColumnAssessment.AssessmentAreaNumber).First(&assessmentArea).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch assessment area"})
			return
		}

		response := responses.ReportColumnAssessmentResponse{
			Number:           reportColumnAssessment.Number,
			AssessmentArea:   reportColumnAssessment.AssessmentAreaNumber,
			AssessmentMetric: reportColumnAssessment.AssessmentMetricNumber,
			Code:             reportColumnAssessment.Code,
			ReportColumn:     reportColumnAssessment.ReportColumnNumber,
		}
		rca = append(rca, response)

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
		Code:                payload.Code,
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
		Code:                 payload.Code,
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
		Code:                 payload.Code,
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

func FetchReportSummaryByStudentAndStudyPeriodAndClass(ctx *gin.Context) {
	studyPeriodNumber := ctx.Param("study_period_number")
	campusClassNumber := ctx.Param("campus_class_number")

	var response []responses.SummarizedFinalReportResponse

	var enrollments []enrollmentModel.Enrollment

	if err := config.DB.Preload("Student").Preload("CampusClass").Where("study_period_number = ? AND campus_class_number = ?", studyPeriodNumber, campusClassNumber).Find(&enrollments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch enrollments"})
		return
	}

	for _, enrollment := range enrollments {
		var reportSummary models.ReportSummary
		if err := config.DB.Where("student_number = ? AND study_period_number = ? AND campus_class_number = ?", enrollment.StudentNumber, enrollment.StudyPeriodNumber, enrollment.CampusClassNumber).First(&reportSummary).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response = append(response, responses.SummarizedFinalReportResponse{
					Student:           enrollment.Student.Surname + " " + enrollment.Student.SecondName,
					Grade:             "N/A",
					Aggregate:         "N/A",
					Class:             enrollment.CampusClass.ClassName,
					StudentNumber:     enrollment.Student.Number,
					EnrollmentNumber:  enrollment.Number,
					CampusClassNumber: enrollment.CampusClassNumber,
				})

			}

		} else {

			response = append(response, responses.SummarizedFinalReportResponse{
				Student:           enrollment.Student.Surname + " " + enrollment.Student.SecondName,
				Grade:             reportSummary.Grade,
				Aggregate:         strconv.Itoa(int(reportSummary.Aggregate)),
				Class:             enrollment.CampusClass.ClassName,
				StudentNumber:     enrollment.Student.Number,
				EnrollmentNumber:  enrollment.Number,
				CampusClassNumber: enrollment.CampusClassNumber,
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func GenerateReportSummary(ctx *gin.Context) {
	var payload requests.GenerateReportSummaryRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err := payload.ValidateGenerateReportSummary()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	//check if report settings exists
	var reportSettings models.ReportSetting
	if err := config.DB.Where("study_period_number,campus_class_number = ? AND status = ?", payload.StudyPeriodNumber, payload.CampusClassNumber, constants.STATUS_ACTIVE).First(&reportSettings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Report settings for this class does not exist"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch settings  for this class"})
		return
	}

	//check if report columns exists
	var reportColumns []models.ReportColumn
	if err := config.DB.Where("report_setting_number = ?", reportSettings.Number).Find(&reportColumns).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Report columns for this class not yet set up"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch columns  for this class"})
		return
	}

	if len(reportColumns) != int(reportSettings.Columns) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Report columns set for the settings are not yet fully set up"})
		return

	}

	reportSummary := models.ReportSummary{
		Number:              uuid.NewV4().String(),
		ReportSettingNumber: reportSettings.Number,
		StudentNumber:       payload.StudentNumber,
		CampusClassNumber:   payload.CampusClassNumber,
		StudyPeriodNumber:   payload.StudyPeriodNumber,
		Grade:               "",
		Aggregate:           -1,
		Comment:             "",
		UserNumber:          ctx.Request.Header["User-Number"][0],
		Code:                0,
	}

	if err := config.DB.Create(&reportSummary).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create report summary"})
		return
	}

	var campusClass classModel.CampusClass

	if err := config.DB.Where("number = ?", payload.CampusClassNumber).First(&campusClass).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch class"})
		return

	}

	var assessmentAreas []assessmentModel.AssessmentArea

	if err := config.DB.Where("curriculum_class_number = ?", campusClass.CurriculumClassNumber).Find(&assessmentAreas).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch assessment areas"})
		return
	}

	for _, assessmentArea := range assessmentAreas {

		for _, reportColum := range reportColumns {

			var rcas []models.ReportColumnAssessment

			if err := config.DB.Where("assessment_area_number = ? AND report_column_number = ? ", assessmentArea.Number, reportColum.Number).Find(&rcas).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					continue
				}
			}

			var score float64
			for _, rca := range rcas {
				var assessment assesmentModel.Assessment

				if err := config.DB.Where("assessment_metric_number = ? AND student_number = ? ", rca.AssessmentMetricNumber, payload.StudentNumber).First(&assessment).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						continue
					}
				}
				score += assessment.ActualScore

			}

			var scoreAverage = score / float64(len(rcas))

			// find the aggregate scale with from grading standard of the  assessment metric

			percentageScore := (scoreAverage / 100) * reportColum.PercentageShare

			var reportScore = models.ReportScore{
				ReportSummaryNumber:  reportSummary.Number,
				AssessmentAreaNumber: assessmentArea.Number,
				ReportColumnNumber:   reportColum.Number,
				Score:                percentageScore,
				AverageScore:         scoreAverage,
				Aggregate:            0,
				Grade:                "",
			}

			config.DB.Create(&reportScore)

		}

		var rs []models.ReportScore

		if err := config.DB.Where("report_summary_number = ? AND assessment_area_number = ?", reportSummary.Number, assessmentArea.Number).Find(&rs).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
		}

		// sum report scores

		var summation float64
		for _, r := range rs {
			summation += r.Score
		}

		var reportScoreSummary = models.ReportScoreSummary{
			ReportSummaryNumber:  reportSummary.Number,
			AssessmentAreaNumber: assessmentArea.Number,
			Aggregate:            0,
			Grade:                "",
			Summation:            summation,
			Comment:              "Subject teacher comment goes here",
		}

		config.DB.Create(&reportScoreSummary)

	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Report summary created successfully"})
}
