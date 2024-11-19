package models

import (
	"time"

	"gorm.io/gorm"
)

type Entity struct {
	gorm.Model
	Number    string `json:"number" gorm:"size:36;unique;not null"`
	Name      string `json:"name" gorm:"size:100;not null"`
	Email     string `json:"email" gorm:"size:100"`
	Phone     string `json:"phone" gorm:"size:100"`
	Address   string `json:"address" gorm:"size:100"`
	AddedBy   string `json:"added_by" gorm:"size:36;not null"`
	Facebook  string `json:"facebook" gorm:"size:100"`
	Instagram string `json:"instagram" gorm:"size:100"`
	Twitter   string `json:"twitter" gorm:"size:100"`
	Website   string `json:"website" gorm:"size:100"`
	Whatsapp  string `json:"whatsapp" gorm:"size:100"`
	Linkedin  string `json:"linkedin" gorm:"size:100"`
}

type GeneralCalendar struct {
	gorm.Model
	EntityNumber   string          `json:"entity_number" gorm:"size:36"`
	Number         string          `json:"number" gorm:"size:36;unique;not null"`
	StartDate      time.Time       `json:"start_date" gorm:"not null"`
	EndDate        time.Time       `json:"end_date" gorm:"not null"`
	Entity         Entity          `gorm:"foreignKey:EntityNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Name           string          `json:"name" gorm:"size:250;not null"`
	Description    string          `json:"description" gorm:"size:250"`
	UserNumber     string          `json:"user_number" gorm:"size:36;not null"`
	CalendarEvents []CalendarEvent `gorm:"foreignKey:CalendarNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type CalendarEvent struct {
	gorm.Model
	Number          string          `json:"number" gorm:"size:36;unique;not null"`
	Name            string          `json:"name" gorm:"size:100;not null"`
	CalendarNumber  string          `json:"calendar_number" gorm:"size:36;not null"`
	StartDate       time.Time       `json:"start_date" gorm:"not null"`
	EndDate         time.Time       `json:"end_date" gorm:"not null"`
	StartTime       string          `json:"start_time" gorm:"size:10;not null"`
	EndTime         string          `json:"end_time" gorm:"size:10;not null"`
	GeneralCalendar GeneralCalendar `gorm:"foreignKey:CalendarNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
