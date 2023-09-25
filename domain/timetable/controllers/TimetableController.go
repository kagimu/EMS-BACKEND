package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/domain/timetable/api/requests"
	"github.com/hamsoftug/wholesomeEd/domain/timetable/models"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func addPeriod(context *gin.Context) {
	var timetable models.Timetable
	var timetableRequest requests.TimetableRequest

	err := context.BindJSON(&timetableRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timetable.Number = uuid.NewV4().String()

}
