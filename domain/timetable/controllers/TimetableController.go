package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/requests"
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
	}

	config.DB.Save(&period)

}

func getDistinctDays(ctx *gin.Context) {
	var days []string
	config.DB.Model(&models.Timetable{}).Select("DISTINCT day").Scan(&days)
	ctx.JSON(http.StatusOK, gin.H{"days": days})

}

func getPeriodsByDay(ctx *gin.Context) {
	day := ctx.Param("day")
	var periods []models.Timetable
	config.DB.Where("day = ?", day).Find(&periods)
	ctx.JSON(http.StatusOK, gin.H{"periods": periods})
}
