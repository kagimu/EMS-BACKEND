package controllers

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assesmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/services"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	learningAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	enrollmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	services3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/services"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"

	"net/http"
)

func CreateReportSettings(c *gin.Context) {

	currentUser, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ReportSettingsRequest
	if err = c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err = payload.ValidateReportSetting()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services3.SaveReportSettings(payload, currentUser.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report settings created successfully"})

}

func FetchReportSettings(c *gin.Context) {

	classNumber := c.Query("class_number")
	campusNumber := c.Query("campus_number")
	campusClassNumber := c.Query("campus_class_number")
	academicYearNumber := c.Query("academic_year_number")
	studyPeriodNumber := c.Query("study_period_number")

	if campusClassNumber == "" && (campusNumber == "" || classNumber == "") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please select at-least a campus and class to view report settings"})
		return
	}

	/*if campusNumber != "" && classNumber != "" {
		campusClassNumber = ""
	}*/

	/*var reportSettings []models.ReportSetting
	if err := config.DB.
		Joins("JOIN campus_classes ON campus_classes.number = assessment_metrics.campus_class_number").
		Joins("JOIN curriculum_classes ON curriculum_classes.number = campus_classes.curriculum_class_number").
		Preload("AcademicYear").
		Preload("CampusClass").
		Preload("StudyPeriod").
		Scopes(
			filters.ByAnyField("campus_classes.campus_number", campusNumber),
			filters.ByAnyField("curriculum_classes.class_number", classNumber),
			filters.ByCampusClassNumber(campusClassNumber),
			filters.ByAcademicYearNumber(c.Query("academic_year_number")),
			filters.ByStudyPeriodNumber(c.Query("study_period_number")),
		).
		Distinct("report_settings.*").
		Find(&reportSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report settings error: " + err.Error()})
		return
	}

	var rs []responses.ReportSettingsResponse

	for _, reportSetting := range reportSettings {

		studyPeriodName, _ := repositories.GetStudyPeriodNameByNumber(reportSetting.StudyPeriod.StudyPeriodNameNumber)

		response := responses.ReportSettingsResponse{
			Name:         reportSetting.AcademicYear.AcademicYearName + " - " + reportSetting.CampusClass.ClassName + " - " + studyPeriodName.Name,
			Number:       reportSetting.Number,
			Status:       reportSetting.Status,
			Columns:      reportSetting.Columns,
			Code:         reportSetting.Code,
			AcademicYear: reportSetting.AcademicYear.AcademicYearName,
			StudyPeriod:  studyPeriodName.Name,
			Class:        reportSetting.CampusClass.ClassName,
		}

		rs = append(rs, response)
	}*/

	allReportSettings, err := services3.FetchAllReportSettings(campusNumber, classNumber, campusClassNumber, academicYearNumber, studyPeriodNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": allReportSettings})

}

func CreatReportColumn(ctx *gin.Context) {

	currentUser, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ReportColumnRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err = payload.ValidateReportColumn()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services3.SaveReportColumn(payload, currentUser.Number)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Report column created successfully"})

}

func FetchReportColumns(ctx *gin.Context) {

	reportSettingsNumber := ctx.Query("report_setting_number")
	if reportSettingsNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please select at least one Report setting"})
		return
	}

	response, err := services3.GetReportColumns(reportSettingsNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})

}

func CreateReportColumnAssessment(ctx *gin.Context) {

	currentUser, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ReportColumnAssessmentRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err = payload.ValidateReportColumnAssessment()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services3.SaveAssessmentToReportColumn(payload, currentUser.Number)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Report column assessment added successfully"})
}

func FetchReportColumnAssessments(ctx *gin.Context) {

	classNumber := ctx.Query("class_number")
	campusNumber := ctx.Query("campus_number")
	campusClassNumber := ctx.Query("campus_class_number")
	curriculumNumber := ctx.Query("curriculum_umber")

	if classNumber != "" && campusNumber != "" {
		campusClass, err := services.ResolveCampusClass(classNumber, campusNumber, curriculumNumber)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"message": err.Error()})
			return
		}

		campusClassNumber = campusClass.Number
	}

	var reportColumnAssessmentResponses = make([]responses.ReportColumnAssessmentResponse, 0)

	var reportColumnAssessments []models.ReportColumnAssessment
	if err := config.DB.Joins("JOIN report_columns ON report_columns.number = report_column_assessments.report_column_number").
		Joins("JOIN report_settings ON report_settings.number = report_columns.report_setting_number").
		Scopes(
			filters.ByAssessmentAreaNumber(ctx.Query("assessment_area_number")),
			filters.ByAssessmentMetricNumber(ctx.Query("assessment_metric_number")),
			filters.ByAnyField("report_column_number", ctx.Query("report_column_number")),
			filters.ByCampusClassNumber(campusClassNumber),
			filters.ByAnyField("report_settings.number", ctx.Query("report_setting_number")),
			filters.ByAnyField("report_settings.academic_year_number", ctx.Query("academic_year_number")),
			filters.ByAnyField("report_settings.study_period_number", ctx.Query("study_period_number")),
		).
		Distinct("report_column_assessments.*").
		Find(&reportColumnAssessments).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, gin.H{"data": reportColumnAssessmentResponses})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch report column assessments error: " + err.Error()})
		return
	}

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
			AssessmentMetric: "Assessment Metric Name",
			ReportColumn:     reportColumn.Name,
		}
		reportColumnAssessmentResponses = append(reportColumnAssessmentResponses, response)

	}

	ctx.JSON(http.StatusOK, gin.H{"data": reportColumnAssessmentResponses})

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

		studyPeriodName, _ := repositories.GetStudyPeriodNameByNumber(reportSummary.StudyPeriod.StudyPeriodNameNumber)

		cleanReportSummary = append(cleanReportSummary, responses.LeanReportSummary{
			Number:        reportSummary.Number,
			StudentName:   reportSummary.Student.Surname + " " + reportSummary.Student.SecondName + " " + reportSummary.Student.OtherName,
			StudentNumber: reportSummary.StudentNumber,
			Grade:         reportSummary.Grade,
			Stream:        reportSummary.Stream.Name,
			StudyPeriod:   studyPeriodName.Name,
			Class:         reportSummary.CampusClass.ClassName,
		})
	}
	paginatedResponse.Content = cleanReportSummary

	ctx.JSON(http.StatusOK, paginatedResponse)

}

func FetchReportSummaryDetails(ctx *gin.Context) {
	reportSummaryNumber := ctx.Param("report_summary_number")
	//htmlReportTemplate := ctx.DefaultQuery("html_report_template", "victorious-report-primary")

	var reportSummary models.ReportSummary
	if err := config.DB.Where("number = ?", reportSummaryNumber).Preload("Student").Preload("StudyPeriod").First(&reportSummary).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report summary"})
		return
	}

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
		assessments := make([]responses.Assessments, 0)

		var assessmentAreas []assessmentAreaModel.AssessmentArea
		if err := config.DB.Where("assessment_areas.curriculum_class_number = ?", campusClass.CurriculumClassNumber).
			Joins("JOIN subjects on subjects.number = assessment_areas.subject_number AND subjects.curriculum_number=?", campusClass.CurriculumClass.CurriculumLevel.CurriculumNumber).
			Joins("JOIN subject_clusters on subject_clusters.number = subjects.subject_cluster_number AND subject_clusters.learning_area_number =?", learningArea.Number).
			Find(&assessmentAreas).Error; err != nil {
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
	scalesUsed, err := reporsitory.GetAppliedScaleByScaleNumber(reportSummary.AggregateScaleUsed)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch scales used"})
		return

	}
	scales := reporsitory.BuildAppliedScaleResponse(scalesUsed)

	studyPeriodName, _ := repositories.GetStudyPeriodNameByNumber(reportSummary.StudyPeriod.StudyPeriodNameNumber)

	response := responses.FullReportSummary{
		Student: responses.LeanStudent{
			Number:     reportSummary.Student.Number,
			Surname:    reportSummary.Student.Surname,
			SecondName: reportSummary.Student.SecondName,
			OtherName:  reportSummary.Student.OtherName,
			Gender:     reportSummary.Student.Gender,
		},
		Grade:               reportSummary.Grade,
		LearningAreas:       la,
		Stream:              "",
		Class:               campusClass.ClassName,
		ClassAbbreviation:   campusClass.ClassAbbreviation,
		School:              campusClass.Campus.School,
		Campus:              campusClass.Campus,
		StudyPeriod:         studyPeriodName.Name,
		ReportColumns:       rcs,
		ReportTemplateName:  reportSummary.ReportTemplateName,
		TeacherComment:      reportSummary.TeacherComment,
		HeadTeachersComment: reportSummary.HeadTeacherComment,
		AggregateScale:      scales,
	}
	ctx.JSON(http.StatusOK, gin.H{"data": response})

	/*//pdfObject, err := pdfGenerator.GeneratePdf(htmlReportTemplate, response)
	pdfObject, err := pdfGenerator.GeneratePdf(htmlReportTemplate, nil)
	if err != nil {
		ctx.JSON(500, gin.H{"message": "Failed to generate pdf" + err.Error()})
		return
	}

	//err = pdfObject.WriteFile("output.pdf")
	ctx.JSON(http.StatusOK, gin.H{"data": base64.StdEncoding.EncodeToString(pdfObject.Bytes())})*/
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
	classNumber := ctx.Query("class_number")
	campusNumber := ctx.Query("campus_number")
	curriculumNumber := ctx.Query("curriculum_number")

	if classNumber != "" && campusNumber != "" {
		campusClass, err := services.ResolveCampusClass(classNumber, campusNumber, curriculumNumber)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if campusClass.Number == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "The selected class is not active on the selected campus"})
			return
		}

		campusClassNumber = campusClass.Number
	}

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

	if err := services3.GenerateReportSummary(payload, ctx.Request.Header["User-Number"][0]); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return

	}

	ctx.JSON(http.StatusOK, commons.Feedback("Report summary generated successfully"))

}

func CreateHeadTeacherComment(ctx *gin.Context) {
	var payload requests.HeadTeacherCommentRequest
	reportSettingNumber := ctx.Param("report_setting_number")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err := payload.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services3.CreateHeadTeacherComment(payload, reportSettingNumber, ctx.Request.Header["User-Number"][0])

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Head teacher comment created successfully"})

}

func FetchHeadTeacherComments(ctx *gin.Context) {
	comments, err := services3.FetchHeadTeacherComments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return

	}
	ctx.JSON(http.StatusOK, comments)

}

func CreateTeachersComment(ctx *gin.Context) {

	var payload requests.TeacherCommentRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err := payload.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var reportSummary models.ReportSummary
	if err := config.DB.Where("number = ?", payload.ReportSummaryNumber).First(&reportSummary).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch given  report summary"})
		return
	}

	// update the teacher_comment field in the report summary
	reportSummary.Comment = payload.Comment
	reportSummary.UserNumber = ctx.Request.Header["User-Number"][0]

	if err := config.DB.Save(&reportSummary).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update teacher comment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Teacher comment created successfully"})

}

func FetchNarrativeScore(c *gin.Context) {

	studentNumber := c.Query("student_number")
	if studentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Provide student number"})
		return
	}

	narrations, err := services2.FetchLearnerNarrative(studentNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": narrations})

}
