package controllers

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	aggregateSaceltModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	assesmentGradeModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	assesmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	learningAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	enrollmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	reportTemplateModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/services"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
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
	if err := config.DB.Where("study_period_number = ? AND academic_year_number = ? AND campus_class_number = ? AND status = ?", payload.StudyPeriodNumber, payload.AcademicYearNumber, payload.CampusClassNumber, constants.StatusActive).First(&existingReportSettings).Error; err == nil {
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
		Status:             constants.StatusActive,
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

	campusClassNumber := c.Query("campus_class_number")

	if campusClassNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please select at-least a class"})
		return
	}

	var reportSettings []models.ReportSetting
	if err := config.DB.Preload("AcademicYear").
		Preload("CampusClass").
		Preload("StudyPeriod").
		Scopes(
			filters.ByCampusClassNumber(campusClassNumber),
			filters.ByAcademicYearNumber(c.Query("academic_year_number")),
			filters.ByStudyPeriodNumber(c.Query("study_period_number")),
		).
		Find(&reportSettings).Error; err != nil {
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

	// validate the current number of entries in the report column table for the report setting

	var reportSetting models.ReportSetting
	if err := config.DB.Where("number = ?", payload.ReportSettingNumber).First(&reportSetting).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report settings"})
		return
	}

	var reportColumns []models.ReportColumn
	if err := config.DB.Where("report_setting_number = ?", payload.ReportSettingNumber).Find(&reportColumns).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report columns"})
		return
	}

	if len(reportColumns) >= int(reportSetting.Columns) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Report columns for this class are already set up"})
		return
	}

	//validate name and code to be unique on each report setting

	var existingReportColumn models.ReportColumn
	if err := config.DB.Where("name = ? AND report_setting_number = ?", payload.Name, payload.ReportSettingNumber).First(&existingReportColumn).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Report column with the same name already exists"})
		return
	}

	if err := config.DB.Where("code = ? AND report_setting_number = ?", payload.Code, payload.ReportSettingNumber).First(&existingReportColumn).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Report column with the same code already exists"})
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
			Name:              reportColumn.Name,
			Number:            reportColumn.Number,
			PercentageShare:   reportColumn.PercentageShare,
			Code:              reportColumn.Code,
			ReportSettingName: reportColumn.ReportSettingName,
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

	//check if report column exists
	var reportColumn models.ReportColumn
	if err := config.DB.Where("number = ?", payload.ReportColumnNumber).First(&reportColumn).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Report column does not exist"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch report column"})
		return

	}

	// check code uniqueness

	var existingReportColumnAssessment models.ReportColumnAssessment
	if err := config.DB.Where("code = ? AND report_column_number = ?", payload.Code, payload.ReportColumnNumber).First(&existingReportColumnAssessment).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Report column assessment with the same code already exists"})
		return
	}

	var existingReportColumn models.ReportColumnAssessment
	if err := config.DB.Where("report_column_number = ? AND assessment_area_number = ? AND assessment_metric_number = ?", payload.ReportColumnNumber, payload.AssessmentAreaNumber, payload.AssessmentMetricNumber).First(&existingReportColumn).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Report column already configured on an  assessment"})
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

		var assessmentMetric assesmentModel.AssessmentMetric
		if err := config.DB.Where("number = ?", reportColumnAssessment.AssessmentMetricNumber).First(&assessmentMetric).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch assessment metric"})
			return
		}

		var reportColumn models.ReportColumn
		if err := config.DB.Where("number = ?", reportColumnAssessment.ReportColumnNumber).First(&reportColumn).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report column"})
			return
		}

		response := responses.ReportColumnAssessmentResponse{
			Number:           reportColumnAssessment.Number,
			AssessmentArea:   assessmentArea.Name,
			AssessmentMetric: assessmentMetric.Name,
			Code:             reportColumnAssessment.Code,
			ReportColumn:     reportColumn.Name,
		}
		rca = append(rca, response)

	}

	ctx.JSON(http.StatusOK, gin.H{"data": rca})

}

func FetchReportSummaries(ctx *gin.Context) {

	var cleanReportSummary = make([]responses.LeanReportSummary, 0)
	var paginatedResponse = pagination.PaginatedResponse{}

	var reportSummaries []models.ReportSummary
	if err := config.DB.Scopes(
		filters.ByStudyPeriodNumber(ctx.Query("study_period_number")),
		filters.ByStudentNumber(ctx.Query("student_number")),
		filters.ByStreamNumber(ctx.Query("stream_number")),
		filters.ByGrade(ctx.Query("grade")),
		filters.ByCampusClassNumber(ctx.Query("campus_class_number")),
		pagination.Paginate(ctx, &paginatedResponse),
	).Preload("Student").Preload("StudyPeriod").
		Preload("Stream").Preload("CampusClass").Find(&reportSummaries).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report summaries"})
		return
	}

	for _, reportSummary := range reportSummaries {
		cleanReportSummary = append(cleanReportSummary, responses.LeanReportSummary{
			Number:        reportSummary.Number,
			StudentName:   reportSummary.Student.Surname + " " + reportSummary.Student.SecondName + " " + reportSummary.Student.OtherName,
			StudentNumber: reportSummary.StudentNumber,
			Grade:         reportSummary.Grade,
			Stream:        reportSummary.Stream.Name,
			StudyPeriod:   reportSummary.StudyPeriod.StudyPeriodName,
			Class:         reportSummary.CampusClass.ClassName,
		})
	}
	paginatedResponse.Content = cleanReportSummary

	ctx.JSON(http.StatusOK, paginatedResponse)

}

func FetchReportSummaryDetails(ctx *gin.Context) {
	reportSummaryNumber := ctx.Param("report_summary_number")

	var reportSummary models.ReportSummary
	if err := config.DB.Where("number = ?", reportSummaryNumber).Preload("Student").Preload("StudyPeriod").First(&reportSummary).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report summary"})
		return
	}

	// get curricullum class number from the campus class number

	var campusClass classModel.CampusClass
	if err := config.DB.Preload("CurriculumClass.CurriculumLevel").Preload("Campus").Preload("Campus.School").Where("number = ?", reportSummary.CampusClassNumber).First(&campusClass).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch campus class"})
		return
	}

	var learningAreas []learningAreaModel.LearningArea
	if err := config.DB.Where("curriculum_number = ?", campusClass.CurriculumClass.CurriculumLevel.CurriculumNumber).Find(&learningAreas).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch learning areas"})
		return
	}

	// find the latest report settings for the class
	var reportSetting models.ReportSetting
	if err := config.DB.Where("campus_class_number = ? AND status = ?", campusClass.Number, constants.StatusActive).First(&reportSetting).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report settings"})
		return

	}

	var rcs []string

	var reportColumns []models.ReportColumn
	if err := config.DB.Where("report_setting_number = ?", reportSetting.Number).Find(&reportColumns).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report columns"})
		return

	}

	for _, reportColumn := range reportColumns {
		rcs = append(rcs, reportColumn.Name)
	}

	var la []responses.LearningArea
	for _, learningArea := range learningAreas {
		var assessments []responses.Assessments

		var assessmentAreas []assessmentAreaModel.AssessmentArea
		if err := config.DB.Where("assessment_areas.curriculum_class_number = ?", campusClass.CurriculumClassNumber).
			Joins("JOIN subjects on subjects.number = assessment_areas.subject_number AND subjects.curriculum_number=?", campusClass.CurriculumClass.CurriculumLevel.CurriculumNumber).Find(&assessmentAreas).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch assessment areas"})
			return
		}

		for _, assessmentArea := range assessmentAreas {

			var reportScoreSummary models.ReportScoreSummary
			if err := config.DB.Where("report_summary_number = ? AND assessment_area_number =? ", reportSummaryNumber, assessmentArea.Number).First(&reportScoreSummary).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report score summaries"})
				return
			}

			var reportScores []models.ReportScore
			if err := config.DB.Where("report_summary_number = ? AND assessment_area_number = ? ", reportSummaryNumber, assessmentArea.Number).Find(&reportScores).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report scores"})
				return
			}

			if len(reportScores) < len(rcs) {
				//  get the difference and add dummy report scores to the report scores
				var diff = len(rcs) - len(reportScores)
				for i := 0; i < diff; i++ {
					reportScores = append(reportScores, models.ReportScore{
						Score: -1,
					})
				}

			}

			assessments = append(assessments, responses.Assessments{
				Name:               assessmentArea.Name,
				ReportScores:       reportScores,
				ReportScoreSummary: reportScoreSummary,
			})
		}

		la = append(la, responses.LearningArea{
			Name:        learningArea.Name,
			Assessments: assessments,
		})

	}

	var school schoolModel.School
	if err := config.DB.Where("number = ?", campusClass.Campus.SchoolNumber).First(&school).Error; err != nil {

	}

	response := responses.FullReportSummary{
		Student: responses.LeanStudent{
			Number:     reportSummary.Student.Number,
			Surname:    reportSummary.Student.Surname,
			SecondName: reportSummary.Student.SecondName,
			OtherName:  reportSummary.Student.OtherName,
			Gender:     reportSummary.Student.Gender,
		},
		Grade:              reportSummary.Grade,
		LearningAreas:      la,
		Stream:             "",
		Class:              campusClass.ClassName,
		ClassAbbreviation:  campusClass.ClassAbbreviation,
		School:             campusClass.Campus.School,
		Campus:             campusClass.Campus,
		StudyPeriod:        reportSummary.StudyPeriod.StudyPeriodName,
		ReportColumns:      rcs,
		ReportTemplateName: reportSummary.ReportTemplateName,
	}
	ctx.JSON(http.StatusOK, gin.H{"data": response})
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

func FetchReportScoreSummaries(ctx *gin.Context) {
	var reportScoreSummaries []models.ReportScoreSummary
	if err := config.DB.Find(&reportScoreSummaries).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report score summaries"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"report_score_summaries": reportScoreSummaries})

}

func FetchReportSummaryByStudentAndStudyPeriodAndClass(ctx *gin.Context) {
	studyPeriodNumber := ctx.Query("study_period_number")
	campusClassNumber := ctx.Query("campus_class_number")

	if studyPeriodNumber == "" || campusClassNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Study period number and campus class number are required"})
		return

	}

	var response []responses.SummarizedFinalReportResponse

	var enrollments []enrollmentModel.Enrollment

	if err := config.DB.
		Preload("Student").Scopes(filters.ByStudyPeriodNumber(studyPeriodNumber), filters.ByCampusClassNumber(campusClassNumber)).
		Find(&enrollments).Error; err != nil {
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
	if err := config.DB.Where("campus_class_number = ? AND status = ?", payload.CampusClassNumber, constants.StatusActive).First(&reportSettings).Error; err != nil {
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
	user := ctx.Request.Header["User-Number"][0]

	tx := config.DB.Begin()
	var stream classModel.Stream
	if err := tx.Where("id > ?", 0).First(&stream).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch stream"})
		return

	}

	var reportTemplate reportTemplateModel.ReportTemplate

	if reportTemplate, err = services.DetermineReportTemplateForClass(payload.CampusClassNumber); err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if reportTemplate.Number == "" {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "This class has no report template set up"})
		return
	}

	// find it there is an existing report summary for the student

	var reportSummary models.ReportSummary
	if err := tx.Where("student_number = ? AND study_period_number = ? AND campus_class_number = ?", payload.StudentNumber, payload.StudyPeriodNumber, payload.CampusClassNumber).First(&reportSummary).Error; err == nil {
		// update the existing report summary with data from the payload
		reportSummary.ReportSettingNumber = reportSettings.Number
		reportSummary.Grade = ""
		reportSummary.Aggregate = -1
		reportSummary.Comment = ""
		reportSummary.UserNumber = user
		reportSummary.StreamNumber = stream.Number
		reportSummary.ReportTemplateName = reportTemplate.Name

		if err := tx.Save(&reportSummary).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update report summary"})
			return
		}

	} else {

		reportSummary = models.ReportSummary{
			Number:              uuid.NewV4().String(),
			ReportSettingNumber: reportSettings.Number,
			StudentNumber:       payload.StudentNumber,
			CampusClassNumber:   payload.CampusClassNumber,
			StudyPeriodNumber:   payload.StudyPeriodNumber,
			Grade:               "",
			Aggregate:           -1,
			Comment:             "",
			UserNumber:          user,
			Code:                0,
			StreamNumber:        stream.Number,
			ReportTemplateName:  reportTemplate.Name,
		}

		if err := tx.Create(&reportSummary).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create report summary"})
			return
		}
	}

	var campusClass classModel.CampusClass

	if err := tx.Where("number = ?", payload.CampusClassNumber).First(&campusClass).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch class"})
		return

	}

	var assessmentAreas []assessmentModel.AssessmentArea

	if err := tx.Where("curriculum_class_number = ?", campusClass.CurriculumClassNumber).Find(&assessmentAreas).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch assessment areas"})
		return
	}

	var totalAggregates float64
	for _, assessmentArea := range assessmentAreas {
		var scale aggregateSaceltModel.AggregateScale

		for _, reportColum := range reportColumns {

			var rcas []models.ReportColumnAssessment

			if err := tx.Where("assessment_area_number = ? AND report_column_number = ? ", assessmentArea.Number, reportColum.Number).Find(&rcas).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					continue
				}
			}

			var score float64
			for _, rca := range rcas {
				var assessment assesmentModel.Assessment

				if err := tx.Where("assessment_metric_number = ? AND student_number = ? ", rca.AssessmentMetricNumber, payload.StudentNumber).First(&assessment).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						continue
					}
				}
				score += assessment.ActualScore

			}
			length := float64(len(rcas))

			var scoreAverage = score / length
			if length == 0 || score == 0 {
				scoreAverage = -1
			}

			// find the aggregate scale with from grading standard of the  assessment metric
			var percentageScore float64

			if length > 0 && score > 0 {
				percentageScore = (scoreAverage / 100) * reportColum.PercentageShare
				scale, err = GetProperAggregateScale(scoreAverage, assessmentArea.Number)
			}

			//check if the report score already exists with the same report summary number, assessment area number and report column number
			var reportScore models.ReportScore
			if err := tx.Where("report_summary_number = ? AND assessment_area_number = ? AND report_column_number = ?", reportSummary.Number, assessmentArea.Number, reportColum.Number).First(&reportScore).Error; err == nil {
				//update the report score
				reportScore.Score = commons.RoundFloat(percentageScore, 2)
				reportScore.AverageScore = commons.RoundFloat(scoreAverage, 2)
				reportScore.Aggregate = scale.Weight
				reportScore.Grade = scale.Name
				reportScore.UserNumber = user

				if score > 0 {
					if err = tx.Save(&reportScore).Error; err != nil {
						tx.Rollback()
						ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update report score"})
						return
					}
				}

			} else {

				//create a new report score

				reportScore = models.ReportScore{
					Number:               uuid.NewV4().String(),
					ReportSummaryNumber:  reportSummary.Number,
					AssessmentAreaNumber: assessmentArea.Number,
					ReportColumnNumber:   reportColum.Number,
					Score:                commons.RoundFloat(percentageScore, 2),
					AverageScore:         commons.RoundFloat(scoreAverage, 2),
					Aggregate:            scale.Weight,
					Grade:                scale.Name,
					UserNumber:           user,
				}

				if score > 0 {
					if err = tx.Create(&reportScore).Error; err != nil {
						tx.Rollback()
						ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create report score"})
						return
					}
				}

			}

		}

		var rs []models.ReportScore

		if err := tx.Where("report_summary_number = ? AND assessment_area_number = ?", reportSummary.Number, assessmentArea.Number).Find(&rs).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
		}
		// sum report scores

		var summation float64
		for _, r := range rs {
			summation += r.Score
		}

		scale, err = GetProperAggregateScale(summation, assessmentArea.Number)

		totalAggregates += scale.Weight

		//check if the report score summary already exists with the same report summary number and assessment area number
		var reportScoreSummary models.ReportScoreSummary
		if err := tx.Where("report_summary_number = ? AND assessment_area_number = ?", reportSummary.Number, assessmentArea.Number).First(&reportScoreSummary).Error; err == nil {
			//update the report score summary
			reportScoreSummary.Aggregate = scale.Weight
			reportScoreSummary.Grade = scale.Name
			reportScoreSummary.Summation = summation
			reportScoreSummary.UserNumber = user

			if err = tx.Save(&reportScoreSummary).Error; err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update report score summary"})
				return
			}

		} else {
			//create a new report score summary
			reportScoreSummary = models.ReportScoreSummary{
				Number:               uuid.NewV4().String(),
				ReportSummaryNumber:  reportSummary.Number,
				AssessmentAreaNumber: assessmentArea.Number,
				Aggregate:            scale.Weight,
				Grade:                scale.Name,
				Summation:            summation,
				Comment:              "Subject teacher comment goes here",
				UserNumber:           user,
			}

			if err = tx.Create(&reportScoreSummary).Error; err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create report score summary"})
				return
			}
		}

	}
	assessmentGrade, err := GetReportSummaryGrade(totalAggregates, payload.CampusClassNumber)

	if err != nil {
		//tx.Rollback()
		//ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get assessment grade"})
		//return
		reportSummary.Grade = "N/A"
	} else {
		reportSummary.Grade = assessmentGrade.Name
	}

	reportSummary.Aggregate = totalAggregates

	if err := tx.Save(&reportSummary).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update report summary"})
		return
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"message": "Report summary created successfully"})

}

func GetProperAggregateScale(score float64, assessmentAreaNumber string) (aggregateSaceltModel.AggregateScale, error) {
	var scale aggregateSaceltModel.AggregateScale
	var aggregateScales []aggregateSaceltModel.AggregateScale
	var aags assessmentAreaModel.AssessmentAreaGradingStandard

	if err := config.DB.Where("assessment_area_number = ?", assessmentAreaNumber).First(&aags).Error; err != nil {
	}
	if err := config.DB.Where("grading_standard_number = ?", aags.GradingStandardNumber).Find(&aggregateScales).Error; err != nil {
	}

	for _, val := range aggregateScales {
		if score >= val.MinValue && score <= val.MaxValue {
			scale = val

			return scale, nil
		}
	}

	return scale, errors.New("no scale found")
}

func GetReportSummaryGrade(points float64, campusClassNumber string) (assesmentGradeModel.AssessmentGrade, error) {
	var assessmentGrade assesmentGradeModel.AssessmentGrade
	var assessmentGrades []assesmentGradeModel.AssessmentGrade

	if err := config.DB.Where("campus_class_number = ?", campusClassNumber).Find(&assessmentGrades).Error; err != nil {
	}

	for _, val := range assessmentGrades {
		if points >= float64(val.MinimumPoints) && points <= float64(val.MaximumPoints) {
			assessmentGrade = val

			return assessmentGrade, nil
		}
	}

	return assessmentGrade, errors.New("no grade found")

}
