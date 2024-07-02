package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/services"
	"net/http"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
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
		ctx.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	err = entityRequest.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	entity.Name = entityRequest.Name
	entity.Phone = entityRequest.Phone
	entity.Email = entityRequest.Email
	entity.Address = entityRequest.Address
	entity.Number = uuid.NewV4().String()
	entity.AddedBy = ctx.GetHeader("User-Number")

	tx := config.DB.Begin()

	if err = tx.Create(&entity).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, commons.Feedback(err.Error()))
		return
	}

	var entityUser models.EntityUser
	entityUser.EntityNumber = entity.Number
	entityUser.UserNumber = ctx.GetHeader("User-Number")
	entityUser.AddedByNumber = ctx.GetHeader("User-Number")

	if err := tx.Create(&entityUser).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, commons.Feedback(err.Error()))
		return

	}

	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, commons.Feedback(err.Error()))
		return

	}

	ctx.JSON(http.StatusOK, commons.Feedback("Entity added successfully"))

}

func AddGeneralCalendar(context *gin.Context) {
	var generalCalendar models.GeneralCalendar
	var generalCalendarRequest requests.GeneralCalendarRequest

	err := context.BindJSON(&generalCalendarRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	// convert start date and end date to time.Time
	startTime, err := time.Parse("2006-01-02", generalCalendarRequest.StartDate)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}
	endDate, err := time.Parse("2006-01-02", generalCalendarRequest.EndDate)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
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
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusCreated, commons.Feedback("Calendar added successfully"))
}

func AddCalendarEvent(ctx *gin.Context) {
	var calendarEvent models.CalendarEvent
	var calendarEventRequest requests.CalendarEventRequest

	err := ctx.BindJSON(&calendarEventRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
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
	calendarEvent.Number = uuid.NewV4().String()

	err = config.DB.Create(&calendarEvent).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commons.Feedback(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, commons.Feedback("Calendar event added successfully"))

}

func GetCalendarDetails(ctx *gin.Context) {

	calendarNumber := ctx.Param("calendar_number")

	var calendar models.GeneralCalendar = models.GeneralCalendar{}

	err := config.DB.Preload("CalendarEvents").Where("number", calendarNumber).First(&calendar).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, calendar)

}

func GetAllCalendars(ctx *gin.Context) {
	var calendars []responses.Calendar
	var generalCalendars []models.GeneralCalendar

	err := config.DB.Find(&generalCalendars).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, calendars)
			return
		}
		ctx.JSON(http.StatusInternalServerError, commons.Feedback(err.Error()))
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

func AddUserToEntity(context *gin.Context) {

	var userEntity models.EntityUser
	var userEntityRequest requests.UserEntityRequest

	err := context.BindJSON(&userEntityRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	userEntity.UserNumber = userEntityRequest.UserNumber
	userEntity.EntityNumber = userEntityRequest.EntityNumber

	err = config.DB.Create(&userEntity).Error
	if err != nil {
		context.JSON(http.StatusInternalServerError, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusOK, commons.Feedback("User added to entity successfully"))
}

func GetUserEntities(ctx *gin.Context) {
	var entities []models.Entity
	user, err := commons.GetUserFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return
	}
	if err := config.DB.Where("entity_users.user_number = ? ", user.Number).Joins("JOIN entity_users ON entity_users.entity_number = entities.number").Find(&entities).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, entities)
}

func SaveCurrentEntity(context *gin.Context) {
	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	entityNumber := context.GetString("entity_number")
	if entityNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Provide entity number"))
		return
	}

	entity, err := services.SaveCurrentEntity(entityNumber, user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": entity, "message": "Current entity set successfully"})
}

func GetCurrentEntity(context *gin.Context) {
	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	entityService := services.EntityService{}

	entity, err := entityService.GetCurrentEntity(user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	if entity.Number == "" {
		context.JSON(http.StatusOK, gin.H{"data": entity, "message": "No active entity set"})
	} else {
		context.JSON(http.StatusOK, gin.H{"data": entity, "message": "Your current active entity"})
	}

}

func UpdateEntity(context *gin.Context) {
	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.UpdateEntityRequest
	if err := context.ShouldBindJSON(&payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	entityService := services.EntityService{}

	err = entityService.UpdateEntity(payload, user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusOK, commons.Feedback("Entity updated successfully"))
}

func DeleteEntity(context *gin.Context) {
	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	entityNumber := context.GetString("entity_number")
	if entityNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Provide entity number to delete"))
		return
	}

	err = services.DeleteEntity(entityNumber, user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusOK, commons.Feedback("Entity removed successfully"))
}

func AllCalendarEvents(context *gin.Context) {
	entityNumber := context.GetString("entity_number")
	if entityNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Provide entity number to get events from."))
		return
	}

	calendars, err := services.FindCalendarEvents(entityNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": calendars})
}
