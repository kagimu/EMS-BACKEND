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
		Scopes(filters.ByDateAfter("calendar_events.end_date", endDateIsAfter)).
		Joins("JOIN general_calendars on general_calendars.number = calendar_events.calendar_number").
		Where("entity_number = ?", entityNumber).
		Order("calendar_events.start_date asc").
		Find(&calendarEvents).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return calendarEvents, nil
		}

		return calendarEvents, err
	}

	return calendarEvents, nil
}

func SaveCalendarEvent(calendarEvent models.CalendarEvent, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Save(&calendarEvent).Error; err != nil {
		return err
	}

	return nil
}

func GetCalendarEventByNumber(calendarEventNumber string, trn ...*gorm.DB) (models.CalendarEvent, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var calendarEvent models.CalendarEvent

	if err := db.Where("number = ?", calendarEventNumber).First(&calendarEvent).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.CalendarEvent{}, nil
		}

		return models.CalendarEvent{}, err
	}

	return calendarEvent, nil
}

func RemoveCalendarEvent(calendarEvent models.CalendarEvent, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Delete(&calendarEvent).Error; err != nil {
		return err
	}

	return nil
}
