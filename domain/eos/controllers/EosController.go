package controllers

import (
	"net/http"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	. "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func AddEntity(ctx *gin.Context) {
	var entity models.Entity
	var entityRequest requests.EntityRequest

	err := ctx.BindJSON(&entityRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Feedback(err.Error()))
		return
	}

	entity.Name = entityRequest.Name
	entity.Phone = entityRequest.Phone
	entity.Email = entityRequest.Email
	entity.Address = entityRequest.Address
	entity.Number = uuid.NewV4().String()

	err = config.DB.Create(&entity).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Feedback(err.Error()))
	}

	ctx.JSON(http.StatusOK, Feedback("Entity added successfully"))

}

func AddGeneralCalendar(context *gin.Context) {
	var generalCalendar models.GeneralCalendar
	var generalCalendarRequest requests.GeneralCalendarRequest

	err := context.BindJSON(&generalCalendarRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, Feedback(err.Error()))
		return
	}

	// convert start date and end date to time.Time
	startTime, err := time.Parse("2006-01-02", generalCalendarRequest.StartDate)
	if err != nil {
		context.JSON(http.StatusBadRequest, Feedback(err.Error()))
		return
	}
	endDate, err := time.Parse("2006-01-02", generalCalendarRequest.EndDate)
	if err != nil {
		context.JSON(http.StatusBadRequest, Feedback(err.Error()))
		return
	}

	generalCalendar.Number = uuid.NewV4().String()
	generalCalendar.Name = generalCalendarRequest.Name
	generalCalendar.StartDate = startTime
	generalCalendar.EndDate = endDate
	generalCalendar.UserNumber = context.GetHeader("User-Number")
	generalCalendar.Description = generalCalendarRequest.Description

	err = config.DB.Create(&generalCalendar).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusCreated, Feedback("Calendar added successfully"))
}

func AddCalendarEvent(ctx *gin.Context) {
	var calendarEvent models.CalendarEvent
	var calendarEventRequest requests.CalendarEventRequest

	err := ctx.BindJSON(&calendarEventRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Feedback(err.Error()))
		return
	}

	//convert start date and end date to time.Time
	startDate, err := time.Parse("2006-01-02", calendarEventRequest.StartDate)
	endDate, err := time.Parse("2006-01-02", calendarEventRequest.EndDate)

	calendarEvent.CalendarNumber = calendarEventRequest.CalendarNumber
	calendarEvent.EndTime = calendarEventRequest.EndTime
	calendarEvent.StartTime = calendarEventRequest.StartTime
	calendarEvent.StartDate = startDate
	calendarEvent.EndDate = endDate

	err = config.DB.Create(&calendarEvent).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Feedback(err.Error()))
		return
	}

}

func GetCalendarDetails(ctx *gin.Context) {

	calendarNumber := ctx.Param("calendar_number")

	var calendar models.GeneralCalendar = models.GeneralCalendar{}

	err := config.DB.Preload("CalendarEvents").Where("number", calendarNumber).First(&calendar).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Feedback(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, calendar)

}

func GetAllCalendars(ctx *gin.Context) {
	var calendars []responses.Calendar = []responses.Calendar{}
	var generalCalendars []models.GeneralCalendar

	err := config.DB.Find(&generalCalendars).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, calendars)
			return
		}
		ctx.JSON(http.StatusInternalServerError, Feedback(err.Error()))
		return
	}

	for _, calendar := range generalCalendars {
		calendars = append(calendars, responses.Calendar{
			Name:         calendar.Name,
			Description:  calendar.Description,
			StartDate:    calendar.StartDate,
			EndDate:      calendar.EndDate,
			Number:       calendar.Number,
			EntityNumber: calendar.EntityNumber,
		})
	}

	ctx.JSON(http.StatusOK, calendars)

}
