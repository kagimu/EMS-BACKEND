package controllers

import (
	"net/http"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	streamModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	subjectModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateStudyingTimeTable(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.StudyingTimeTableRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	var studyingBlockNumber models.StudyingTimeTable
	if err := config.DB.Where("number = ?", payload.StudyingTimeBlockNumber).First(&studyingBlockNumber).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Studying time block not found"})
		return
	}

	var stream streamModels.Stream
	if err := config.DB.Where("number = ?", payload.StreamNumber).First(&stream).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Stream not found"})
		return
	}

	var subject subjectModels.Subject
	if err := config.DB.Where("number = ?", payload.SubjectNumber).First(&subject).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Stream not found"})
		return
	}

	formatStartTime, _ := time.Parse("15:06", payload.StartTime)
	formatEndTime, _ := time.Parse("15:06", payload.EndTime)
	formatDay, _ := time.Parse("Monday", payload.Day)

	saveStudyingTimeTable := models.StudyingTimeTable{
		Number:                  uuid.NewV4().String(),
		StudyingTimeBlockNumber: payload.StudyingTimeBlockNumber,
		StreamNumber:            payload.StreamNumber,
		SubjectNumber:           payload.SubjectNumber,
		Day:                     formatDay,
		StartTime:               formatStartTime,
		EndTime:                 formatEndTime,
		UserNumber:              user.Number,
	}

	var studyingTimeBlock models.StudyingTimeBlock
	if saveStudyingTimeTable.EndTime.After(studyingTimeBlock.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "End time cannot be after the end time of the studying time block"})
		return
	}

	if err := config.DB.Create(&saveStudyingTimeTable).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save studying time table"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Studying time table saved successfully"})
}

func FetchStudyingTimeTables(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var studyingTimeTables []models.StudyingTimeTable
	if err := config.DB.Preload("StudyingTimeBlock").Preload("Stream").Preload("Subject").Find(&studyingTimeTables).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch studying time tables"})
		return
	}

	var response []responses.StudyingTimeTableResponse
	for _, studyingTimeTable := range studyingTimeTables {
		response = append(response, responses.StudyingTimeTableResponse{
			Number:                     studyingTimeTable.Number,
			StudyingTimeBlockNumber:    studyingTimeTable.Number,
			StudyingTimeBlockStartTime: studyingTimeTable.StartTime,
			Day:                        studyingTimeTable.Day,
			StudyingTimeBlockEndTime:   studyingTimeTable.EndTime,
			StreamName:                 studyingTimeTable.Stream.Name,
			SubjectName:                studyingTimeTable.Subject.Name,
			StartTime:                  studyingTimeTable.StartTime,
			EndTime:                    studyingTimeTable.EndTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
