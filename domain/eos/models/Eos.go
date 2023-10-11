package models

import (
	"gorm.io/gorm"
	"time"
)

type GeneralCalendar struct {
	gorm.Model
	Number    string    `json:"number" gorm:"size:36;unique;not null"`
	StartDate time.Time `json:"start_date" gorm:"not null"`
	EndDate   time.Time `json:"end_date" gorm:"not null"`
	Name      string    `json:"name" gorm:"size:36;not null"`
}

type CalendarEvent struct {
	gorm.Model
	Number    string    `json:"number" gorm:"size:36;unique;not null"`
	StartDate time.Time `json:"start_date" gorm:"not null"`
	EndDate   time.Time `json:"end_date" gorm:"not null"`
	StartTime string    `json:"start_time" gorm:"size:10;not null"`
	EndTime   string    `json:"end_time" gorm:"size:10;not null"`
}
