package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	assessmentMetricModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/services"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	enrollmentResponse "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/responses"
	enrollmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var narrativeService = services.NewNarrativeService()

func SaveNarrative(ctx *gin.Context) {

	user, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var request requests.NarrativeRequest
	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = narrativeService.CreateNarrative(request, user.Number)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Narrative saved"})
}

func GetNarrative(context *gin.Context) {

	_startDate := context.Query("start_date")
	_endDate := context.Query("end_date")

	var startDate time.Time
	var endDate time.Time
	var err error

	if _startDate != "" {
		startDate, err = time.Parse(time.DateOnly, _startDate)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Error reading start date: " + err.Error()})
			return
		}
	}

	if _endDate != "" {
		endDate, err = time.Parse(time.DateOnly, _endDate)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Error reading end date: " + err.Error()})
			return
		}
	}

	narrative := requests.NarrativeRequest{
		Description:          context.Query("description"),
		StudentNumber:        context.Query("student_number"),
		AssessmentAreaNumber: context.Query("assessment_area_number"),
		StudyPeriodNumber:    context.Query("study_period_number"),
		CampusClassNumber:    context.Query("campus_class_number"),
	}

	data, err := narrativeService.AllNarrative(narrative, startDate, endDate)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": data})
}
func GetNarrativeStudents(c *gin.Context) {

	assessmentMetricNumber := c.Query("narrative_schedule_number")
	studyPeriodNumber := c.Query("study_period_number")

	if assessmentMetricNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Narrative schedule is required"})
		return
	}

	var assessmentMetric assessmentMetricModel.ScheduledNarrative

	if err := config.DB.Where("number = ?", assessmentMetricNumber).First(&assessmentMetric).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Narrative schedule not found"})
		return
	}

	//fetch enrollment students from assessment metric campus class number

	var enrollments []enrollmentModel.Enrollment

	if err := config.DB.Where("campus_class_number = ? AND status = ? AND study_period_number = ?", assessmentMetric.CampusClassNumber, constants.StatusActive, studyPeriodNumber).Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch students"})
		return
	}

	var students []enrollmentResponse.NarrativeEnrolledStudent

	for _, enrollment := range enrollments {
		var student models.Student
		if err := config.DB.Where("number = ?", enrollment.StudentNumber).First(&student).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch student"})
			return
		}

		// find assesment for student and assessment metric

		var assessment assessmentMetricModel.Narrative

		if err := config.DB.Where("student_number = ? AND narrative_schedule_number = ?", student.Number, assessment.Number).First(&assessment).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch narrative assessments " + err.Error()})
				return
			}
		}

		var campusClass classModels.CampusClass
		if err := config.DB.Where("number = ?", enrollment.CampusClassNumber).First(&campusClass).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Campus class not found"})
			return
		}

		students = append(students, enrollmentResponse.NarrativeEnrolledStudent{
			Number:            student.Number,
			StudentSurname:    student.Surname,
			StudentOtherName:  student.OtherName,
			StudentSecondName: student.SecondName,
			Narrative:         assessment.Description,
			SuggestedSolution: assessment.SuggestedSolution,
			Class:             campusClass.ClassName,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": students})
}

func SaveNarrativeSchedule(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var request requests.ScheduledNarrativeRequest
	err = context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = narrativeService.SaveSchedule(request, user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Scheduled Narrative saved"})
}

func NarrativeSchedule(context *gin.Context) {
	studyPeriodNumber := context.Query("study_period_number")
	classNumber := context.Query("class_number")
	campusNumber := context.Query("campus_number")
	curriculumNumber := context.Query("curriculum_number")
	assessmentAreaNumber := context.Query("assessment_area_number")

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "choose class"})
		return
	}

	if campusNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "choose campus"})
		return
	}

	if curriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "choose curriculum"})
		return
	}

	if studyPeriodNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "study_period_number is required"})
		return
	}

	data, err := narrativeService.AllSchedules(studyPeriodNumber, classNumber, campusNumber, curriculumNumber, assessmentAreaNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": data})
}

func GetDetailedNarrative(context *gin.Context) {
	narrativeNumber := context.Query("narrative_number")
	if narrativeNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "narrative number is required"})
		return
	}

	data, err := narrativeService.NarrativeDetails(narrativeNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": data})
}

func DeleteNarrative(context *gin.Context) {
	narrativeNumber := context.Query("narrative_number")
	if narrativeNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "narrative number is required"})
		return
	}

	err := narrativeService.DeleteNarrative(narrativeNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Removed successfully"})
}

func DeleteNarrativeSchedule(context *gin.Context) {
	scheduleNumber := context.Param("schedule_number")
	if scheduleNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "schedule number is required"})
		return
	}

	err := narrativeService.DeleteNarrativeSchedule(scheduleNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Removed successfully"})
}

func UpdateNarrative(ctx *gin.Context) {
	narrativeNumber := ctx.Param("narrative_number")

	if narrativeNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "narrative number is required"})
		return
	}

	user, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var request requests.NarrativeRequest
	err = ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = narrativeService.UpdateNarrativeEntry(request, narrativeNumber, user.Number)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Entry updated successfully"})
}

func UpdateNarrativeSchedule(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	scheduleNumber := context.Param("schedule_number")
	if scheduleNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "schedule number is required"})
		return
	}

	var request requests.ScheduledNarrativeRequest
	err = context.ShouldBindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = narrativeService.UpdateNarrativeSchedule(request, scheduleNumber, user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Schedule updated successfully"})
}

func NarrativePrintPreview(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	_startDate := context.Query("start_date")
	_endDate := context.Query("end_date")

	var startDate time.Time
	var endDate time.Time

	if _startDate != "" {
		startDate, err = time.Parse(time.DateOnly, _startDate)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Error reading start date: " + err.Error()})
			return
		}
	}

	if _endDate != "" {
		endDate, err = time.Parse(time.DateOnly, _endDate)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Error reading end date: " + err.Error()})
			return
		}
	}

	narrative := requests.NarrativeRequest{
		StudentNumber:     context.Query("student_number"),
		StudyPeriodNumber: context.Query("study_period_number"),
		CampusClassNumber: context.Query("campus_class_number"),
	}

	data, err := narrativeService.NarrativeReportPrintPreview(narrative, startDate, endDate, user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": data})
}

func NarrativeScriptMarkingGuide(c *gin.Context) {

	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	narrativeNumber := c.PostForm("narrative_number")
	if narrativeNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide the assessment schedule"})
		return
	}

	err = services.UploadNarrativeMarkingGuide(c, narrativeNumber, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("Marking guide uploaded successfully"))
}
