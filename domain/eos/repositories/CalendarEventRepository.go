package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"gorm.io/gorm"
	"time"
)

func FindCalendarEvent(entityNumber string, endDateIsAfter time.Time) ([]models.CalendarEvent, error) {
	db := config.DB

	var calendarEvents []models.CalendarEvent

	if err := db.
		Scopes(filters.ByDateAfter("end_date", endDateIsAfter)).
		Joins("JOIN general_calendars on general_calendars.number = calendar_events.calendar_number").
		Where("entity_number = ?", entityNumber).
		Find(&calendarEvents).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return calendarEvents, nil
		}

		return calendarEvents, err
	}

	return calendarEvents, nil
}
