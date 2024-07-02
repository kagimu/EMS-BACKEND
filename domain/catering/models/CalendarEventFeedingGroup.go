package models

import (
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
)

type CalendarEventFeedingGroup struct {
	Number              string               `json:"number" gorm:"size:36;unique;not null"`
	FeedingGroupNumber  string               `json:"feeding_group_number" gorm:"size:36;primaryKey;not null"`
	FeedingGroup        FeedingGroup         `gorm:"foreignKey:FeedingGroupNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	CalendarEventNumber string               `json:"calendar_event_number" gorm:"size:36;primaryKey;not null"`
	CalendarEvent       models.CalendarEvent `gorm:"foreignKey:CalendarEventNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber          string               `json:"user_number" gorm:"size:36;not null"`
	User                userModel.User       `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
