package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func FindCalendarEventFeedingGroupsByGroupNumber(groupNumber string, trn ...*gorm.DB) ([]models.CalendarEventFeedingGroup, error) {
	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	var calendarEventFeedingGroups []models.CalendarEventFeedingGroup

	if err := db.Where("feeding_group_number = ?", groupNumber).Find(&calendarEventFeedingGroups).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return calendarEventFeedingGroups, nil
		}

		return calendarEventFeedingGroups, err
	}

	return calendarEventFeedingGroups, nil
}

func GetCalendarEventFeedingGroupsByGroupNumber(groupNumber string, calendarEventNumber string, trn ...*gorm.DB) (models.CalendarEventFeedingGroup, error) {
	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	var calendarEventFeedingGroups models.CalendarEventFeedingGroup

	if err := db.Where("feeding_group_number = ? AND calendar_event_number = ? ", groupNumber, calendarEventNumber).Find(&calendarEventFeedingGroups).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return calendarEventFeedingGroups, nil
		}

		return calendarEventFeedingGroups, err
	}

	return calendarEventFeedingGroups, nil
}

func RemoveCalendarEventFeedingGroup(calendarEventFeedingGroup models.CalendarEventFeedingGroup, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Delete(&calendarEventFeedingGroup).Error; err != nil {
		return err
	}

	return nil
}

func SavedCalendarEventFeedingGroup(calendarEventFeedingGroup models.CalendarEventFeedingGroup, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Save(&calendarEventFeedingGroup).Error; err != nil {
		return err
	}

	return nil
}
