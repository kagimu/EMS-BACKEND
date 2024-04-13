package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	lessonPlanStepModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateStudyPeriod(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.StudyPeriodRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate payload"})
		return
	}

	var studyPeriodTimeTable models.StudyingTimeTable
	if err := config.DB.Where("number = ?", payload.StudyingTimeTableNumber).First(&studyPeriodTimeTable).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Studying time table not found"})
		return
	}

	var lessonPlanStep lessonPlanStepModels.LessonPlanStep
	if err := config.DB.Where("number = ?", payload.LessonPlanStepNumber).First(&lessonPlanStep).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Lesson plan step not found"})
		return
	}

	saveStudyPeriod := models.StudyPeriod{
		Number:                  uuid.NewV4().String(),
		Comment:                 payload.Comment,
		StudyingTimeTableNumber: payload.StudyingTimeTableNumber,
		UserNumber:              user.Number,
		LessonPlanStepNumber:    payload.LessonPlanStepNumber,
	}

	if err := config.DB.Create(&saveStudyPeriod).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save study period"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Study period saved successfully"})
}

func FetchStudyPeriods(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var studyingPeriods []models.StudyPeriod
	if err := config.DB.Preload("StudyingTimeBlock").Preload("Stream").Preload("Subject").Find(&studyingPeriods).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch studying time tables"})
		return
	}

	var response []responses.StudyPeriodResponse
	for _, studyingPeriod := range studyingPeriods {
		response = append(response, responses.StudyPeriodResponse{
			Number:                     studyingPeriod.Number,
			UserName:                   studyingPeriod.User.FirstName + " " + studyingPeriod.User.SecondName + " " + studyingPeriod.User.LastName,
			LessonPlanStepName:         studyingPeriod.LessonPlanStep.Step,
			StudyingTimeTableNumber:    studyingPeriod.StudyingTimeTableNumber,
			StudyingTimeTableStartTime: studyingPeriod.StudyingTimeTable.StartTime,
			StudyingTimeTableEndTime:   studyingPeriod.StudyingTimeTable.EndTime,
			Comment:                    studyingPeriod.Comment,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
