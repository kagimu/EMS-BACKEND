package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusModes "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	subjectModes "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

func AddPeriod(context *gin.Context) {
	var timetable models.Timetable
	var timetableRequest requests.TimetableRequest

	err := context.BindJSON(&timetableRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timetable.Number = uuid.NewV4().String()
	startTime, err := time.Parse("15:04:05", timetableRequest.StartTime)
	endTime, err := time.Parse("15:04:05", timetableRequest.EndTime)

	period := models.Timetable{
		Number:        timetable.Number,
		StartTime:     startTime,
		EndTime:       endTime,
		SubjectNumber: timetableRequest.SubjectNumber,
		StreamNumber:  timetable.StreamNumber,
		Day:           timetableRequest.Day,
		CampusNumber:  timetableRequest.CampusNumber,
	}

	config.DB.Begin()

	config.DB.Save(&period)

	config.DB.Commit()

	context.JSON(http.StatusOK, gin.H{"message": "Period added successfully"})

}

func GetDistinctDays(ctx *gin.Context) {
	var days []string
	config.DB.Model(&models.Timetable{}).Select("DISTINCT day").Scan(&days)
	ctx.JSON(http.StatusOK, gin.H{"days": days})

}

func GetPeriodsByDay(ctx *gin.Context) {
	day := ctx.Param("day")
	var periods []models.Timetable
	config.DB.Where("day = ?", day).Find(&periods)
	ctx.JSON(http.StatusOK, gin.H{"data": periods})
}

func GetPeriodsByStream(ctx *gin.Context) {
	streamNumber := ctx.Param("stream_number")
	var periods []models.Timetable
	config.DB.Where("stream_number = ?", streamNumber).Find(&periods)
	ctx.JSON(http.StatusOK, gin.H{"data": periods})
}

func FetchClassLevelsByCampusNumberAndLevelNumber(ctx *gin.Context) {

	campusNumber := ctx.Query("campus_number")
	levelNumber := ctx.Query("level_number")

	fmt.Println(campusNumber, levelNumber)

	var classLevels []classModels.ClassLevel

	config.DB.
		Where("campus_number = ? AND level_number = ?", campusNumber, levelNumber).
		Preload("Class").
		Find(&classLevels)
	ctx.JSON(http.StatusOK, gin.H{"data": classLevels})

}

func FetchSubjectsByClassLevelNumber(ctx *gin.Context) {
	classLevelNumber := ctx.Param("class_level_number")
	var subjects []responses.Subjects

	config.DB.Model(subjectModes.Subject{}).
		Select("subjects.name,subjects.number").
		Joins("left join class_level_subjects on class_level_subjects.subject_number=subjects.number").
		Where("class_level_subjects.class_level_number = ?", classLevelNumber).
		Scan(&subjects)
	ctx.JSON(http.StatusOK, gin.H{"data": subjects})
}

func FetchLevelsByCampusNumber(ctx *gin.Context) {
	campusNumber := ctx.Param("campus_number")
	var levels []campusModes.Level

	config.DB.
		Where("campus_number = ?", campusNumber).
		Find(&levels)

	ctx.JSON(http.StatusOK, gin.H{"data": levels})
}
