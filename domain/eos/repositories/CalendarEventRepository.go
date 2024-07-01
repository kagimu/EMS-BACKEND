package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"gorm.io/gorm"
)

func FindCalendarEvent(entityNumber string) ([]models.CalendarEvent, error) {
	db := config.DB

	var calendarEvents []models.CalendarEvent

	if err := db.
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
