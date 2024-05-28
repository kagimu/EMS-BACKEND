package controllers

import (
	"net/http"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	academicYearModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	levelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateStudyingTimeBlock(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.StudyTimeBlockRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	var academicYear academicYearModels.AcademicYear
	if err := config.DB.Where("number = ?", payload.AcademicYearNumber).First(&academicYear).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Academic year not found"})
		return
	}

	var campus campusModels.Campus
	if err := config.DB.Where("number = ?", payload.CampusNumber).First(&campus).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Campus not found"})
		return
	}

	var level levelModels.Level
	if err := config.DB.Where("number = ?", payload.LevelNumber).First(&level).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Level not found"})
		return
	}

	formatStartTime, _ := time.Parse("15:04", payload.StartTime)
	formatEndTime, _ := time.Parse("15:04", payload.EndTime)

	saveStudyingTimeBlock := models.StudyingTimeBlock{
		Number:             uuid.NewV4().String(),
		AcademicYearNumber: payload.AcademicYearNumber,
		CampusNumber:       payload.CampusNumber,
		LevelNumber:        payload.LevelNumber,
		StartTime:          formatStartTime,
		EndTime:            formatEndTime,
		Length:             payload.Length,
		UserNumber:         user.Number,
	}

	if err := config.DB.Create(&saveStudyingTimeBlock).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save studying time block"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Studying time block saved successfully"})
}

func FetchStudyingTimeBlocks(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var studyingTimeBlocks []models.StudyingTimeBlock
	if err := config.DB.Preload("AcademicYear").Preload("Campus").Preload("Level").Find(&studyingTimeBlocks).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch studying time blocks"})
		return
	}

	var response []responses.StudyingTimeBlockResponse
	for _, studyingTimeBlock := range studyingTimeBlocks {
		response = append(response, responses.StudyingTimeBlockResponse{
			Number:             studyingTimeBlock.Number,
			StartTime:          studyingTimeBlock.StartTime,
			EndTime:            studyingTimeBlock.EndTime,
			Length:             uint(studyingTimeBlock.Length),
			AcademicYearNumber: studyingTimeBlock.AcademicYearNumber,
			AcademicYearName:   studyingTimeBlock.AcademicYear.AcademicYearName,
			CampusName:         studyingTimeBlock.Campus.Name,
			LevelName:          studyingTimeBlock.Level.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
