package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	classLevelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
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

	var classLevels []responses.Classes

	config.DB.Model(classLevelModels.ClassLevel{}).
		Select("class_levels.number number,class.name name").
		Joins("left join class on classes.number=class_levels.class_number").
		Where("class_levels.campus_number = ? AND class_levels.level_number = ?", campusNumber, levelNumber).
		Scan(&classLevels)
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
