package controllers

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	enrollmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	services3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/services"
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

func AllReportColumnAssessments(ctx *gin.Context) {

	classNumber := ctx.Query("class_number")
	campusNumber := ctx.Query("campus_number")
	campusClassNumber := ctx.Query("campus_class_number")
	curriculumNumber := ctx.Query("curriculum_number")
	campusCurriculumNumber := ctx.Query("campus_curriculum_number")
	reportColumnNumber := ctx.Query("report_column_number")
	assessmentAreaNumber := ctx.Query("assessment_area_number")
	assessmentMetricNumber := ctx.Query("assessment_metric_number")
	academicYearNumber := ctx.Query("academic_year_number")
	studyPeriodNumber := ctx.Query("study_period_number")
	reportSettingNumber := ctx.Query("report_setting_number")

	if classNumber != "" && campusNumber != "" {
		campusClass, err := services.ResolveCampusClass(classNumber, campusNumber, curriculumNumber, campusCurriculumNumber)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"message": err.Error()})
			return
		}

		campusClassNumber = campusClass.Number
	}

	reportColumnAssessmentResponses, err := services3.AllReportColumnAssessments(reportColumnNumber, assessmentAreaNumber, assessmentMetricNumber, academicYearNumber, studyPeriodNumber, campusClassNumber, reportSettingNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": reportColumnAssessmentResponses})

}

func FetchReportSummaries(ctx *gin.Context) {

	var cleanReportSummary = make([]responses.LeanReportSummary, 0)
	var paginatedResponse = pagination.PaginatedResponse{}

	classNumber := ctx.Query("class_number")
	campusNumber := ctx.Query("campus_number")
	curriculumNumber := ctx.Query("curriculum_number")

	campusClassNumber := ctx.Query("campus_class_number")

	if campusClassNumber == "" {
		campusClass, err := services.ResolveCampusClass(classNumber, campusNumber, curriculumNumber, "")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		campusClassNumber = campusClass.Number
	}

	if campusClassNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please select  a class to view report summaries"})
		return
	}

	var reportSummaries []models.ReportSummary
	if err := config.DB.Scopes(
		filters.ByStudyPeriodNumber(ctx.Query("study_period_number")),
		filters.ByStudentNumber(ctx.Query("student_number")),
		filters.ByStreamNumber(ctx.Query("stream_number")),
		filters.ByGrade(ctx.Query("grade")),
		filters.ByCampusClassNumber(campusClassNumber),
		pagination.Paginate(ctx, &paginatedResponse),
	).Preload("Student").Preload("StudyPeriod").
		Preload("Stream").Preload("CampusClass").Find(&reportSummaries).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch report summaries"})
		return
	}

	for _, reportSummary := range reportSummaries {

		studyPeriodName, _ := repositories.GetStudyPeriodNameByNumber(reportSummary.StudyPeriod.StudyPeriodNameNumber)

		promotionStatus := ""

		promotion, err := repositories2.FindPromotionByReportSummaryNumber(reportSummary.Number)
		if err == nil {
			promotionStatus = promotion.Status
		} else {
			fmt.Println(fmt.Sprintf("Checking promotion status error %s", err.Error()))
		}

		leanStudent := services3.CreateLeanStudent(reportSummary.Student)

		cleanReportSummary = append(cleanReportSummary, responses.LeanReportSummary{
			Number:          reportSummary.Number,
			StudentName:     leanStudent.Surname + " " + leanStudent.SecondName + " " + leanStudent.OtherName,
			StudentNumber:   reportSummary.StudentNumber,
			Photo:           commons.GenerateDownloadPath(reportSummary.Student.Photo, true),
			Grade:           reportSummary.Grade,
			Stream:          reportSummary.Stream.Name,
			StudyPeriod:     studyPeriodName.Name,
			Class:           reportSummary.CampusClass.ClassName,
			PromotionStatus: promotionStatus,
		})
	}
	paginatedResponse.Content = cleanReportSummary

	ctx.JSON(http.StatusOK, paginatedResponse)

}

func FetchReportSummaryDetails(ctx *gin.Context) {
	reportSummaryNumber := ctx.Param("report_summary_number")

	response, err := services3.GetReportSummary(reportSummaryNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
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
	classNumber := ctx.Query("class_number")
	campusNumber := ctx.Query("campus_number")
	curriculumNumber := ctx.Query("curriculum_number")
	campusCurriculumNumber := ctx.Query("campus_curriculum_number")

	if classNumber != "" && campusNumber != "" {
		campusClass, err := services.ResolveCampusClass(classNumber, campusNumber, curriculumNumber, campusCurriculumNumber)
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

		student := services3.CreateLeanStudent(enrollment.Student)

		var reportSummary models.ReportSummary
		if err := config.DB.Where("student_number = ? AND study_period_number = ? AND campus_class_number = ?", enrollment.StudentNumber, enrollment.StudyPeriodNumber, enrollment.CampusClassNumber).First(&reportSummary).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response = append(response, responses.SummarizedFinalReportResponse{
					Student:           student.Surname + " " + student.SecondName,
					Grade:             "N/A",
					Aggregate:         "N/A",
					Class:             enrollment.CampusClass.ClassName,
					StudentNumber:     student.Number,
					EnrollmentNumber:  enrollment.Number,
					CampusClassNumber: enrollment.CampusClassNumber,
				})

			}

		} else {

			response = append(response, responses.SummarizedFinalReportResponse{
				Student:           student.Surname + " " + student.SecondName,
				Grade:             reportSummary.Grade,
				Aggregate:         strconv.Itoa(int(reportSummary.Aggregate)),
				Class:             enrollment.CampusClass.ClassName,
				StudentNumber:     student.Number,
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

func UpdateReportSettings(context *gin.Context) {
	currentUser, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.UpdateReportSettingsRequest
	if err = context.ShouldBindJSON(&payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err = payload.ValidateReportSetting()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services3.UpdateReportSettings(payload, currentUser.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Report settings updated successfully"})
}

func DeleteReportColumns(context *gin.Context) {
	reportColumnNumber := context.Param("report_column_number")
	if reportColumnNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Provide report column number"})
		return
	}

	err := services3.DeleteReportColumn(reportColumnNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Report column deleted successfully"})
}

func DeleteReportSetting(context *gin.Context) {
	reportSettingNumber := context.Query("report_setting_number")
	if reportSettingNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Provide report setting number"})
		return
	}

	err := services3.DeleteReportSetting(reportSettingNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Report setting deleted successfully"})
}

func DeleteReportColumnAssessment(context *gin.Context) {
	reportColumnAssessmentNumber := context.Query("report_column_assessment_number")
	if reportColumnAssessmentNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Provide report column assessment number"})
		return
	}

	err := services3.DeleteReportColumnAssessment(reportColumnAssessmentNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Report column assessment deleted successfully"})
}

func FetchReportPromotionDetails(ctx *gin.Context) {
	//reportSummaryNumber := ctx.Param("report_summary_number")
	//
	//response, err := services3.FetchPromotionStatus(reportSummaryNumber)
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	//}

	ctx.JSON(http.StatusOK, gin.H{"data": ""})
}

func DeleteHeadTeacherComment(ctx *gin.Context) {
	commentNumber := ctx.Param("id")
	if commentNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Provide comment number"})
		return
	}
	commentNumberInt, err := strconv.Atoi(commentNumber)

	err = services3.DeleteHeadTeacherComment(uint(commentNumberInt))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Head teacher comment deleted successfully"})
}

func UpdateReportColumnAssessment(context *gin.Context) {
	reportColumnAssessmentNumber := context.Param("report_column_assessment_number")
	if reportColumnAssessmentNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Provide report column assessment number"})
		return
	}

	currentUser, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.ReportColumnAssessmentUpdateRequest
	if err := context.ShouldBindJSON(&payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err = payload.ValidateReportColumnAssessment()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services3.UpdateReportColumnAssessment(payload, reportColumnAssessmentNumber, currentUser)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Report column assessment updated successfully"})
}

func CreateLearningAreaComment(context *gin.Context) {

	var payload requests.LearningAreaReportCommentRequest

	if err := context.ShouldBindJSON(&payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input body"})
		return
	}

	err := payload.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	currentUser, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	err = services3.SaveLearningArea(payload, currentUser)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Learning area comment updated successfully"})
}

func CreateClassTeacherComment(context *gin.Context) {
	var payload requests.ClassTeacherCommentRequest

	if err := context.ShouldBindJSON(&payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input body" + err.Error()})
		return
	}

	err := payload.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	currentUser, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	err = services3.SaveClassTeacherComment(payload, currentUser)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Class teacher comment created successfully"})
}

func FetchClassTeacherComments(context *gin.Context) {
	//campusClassNumber := context.Query("campus_class_number")
	//if campusClassNumber == "" {
	//	context.JSON(http.StatusBadRequest, gin.H{"message": "Please select a class to view comments"})
	//	return
	//}

	classTeacherComments, err := services3.FetchClassTeacherComments()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": classTeacherComments})
}

func DeleteClassTeacherComment(context *gin.Context) {
	commentNumber := context.Param("id")
	if commentNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Provide comment number"})
		return
	}

	err := services3.DeleteClassTeacherComment(commentNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Class teacher comment deleted successfully"})
}

func CreateSubjectTeacherComment(context *gin.Context) {
	var payload requests.SubjectTeacherCommentRequest

	if err := context.ShouldBindJSON(&payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input body"})
		return
	}

	err := payload.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	currentUser, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	err = services3.SaveSubjectTeacherComment(payload, currentUser)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": errors.New("failed to save subject teacher comment " + err.Error())})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Subject teacher comment created successfully"})
}

func FetchSubjectTeacherComments(context *gin.Context) {
	campusClassNumber := context.Query("campus_class_number")
	if campusClassNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Please select a class to view comments"})
		return
	}

	subjectTeacherComments, err := services3.FetchSubjectTeacherComments()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": subjectTeacherComments})
}

func DeleteSubjectTeacherComment(context *gin.Context) {
	commentNumber := context.Param("id")
	if commentNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Provide comment number"})
		return
	}

	err := services3.DeleteSubjectTeacherComment(commentNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Subject teacher comment deleted successfully"})
}
