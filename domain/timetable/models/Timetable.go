package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	subjectModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"gorm.io/gorm"
	"time"
)

type Timetable struct {
	gorm.Model
	Number        string                `json:"number" gorm:"size:36;unique;not null"`
	StreamNumber  string                `json:"stream_number" gorm:"size:36;not null"`
	Day           string                `json:"day" gorm:"size:36;not null"`
	DayCode       int32                 `json:"day_code" gorm:"size:36;not null"`
	StartTime     string                `json:"start_time" gorm:"size:10;not null"`
	EndTime       string                `json:"end_time" gorm:"size:10;not null"`
	SubjectNumber string                `json:"subject_number" gorm:"size:36;not null"`
	CampusNumber  string                `json:"campus_number" gorm:"size:36;not null"`
	UserNumber    string                `json:"user_number" gorm:"size:36;not null"`
	User          authModels.User       `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Campus        campusModels.Campus   `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Subject       subjectModels.Subject `gorm:"foreignKey:SubjectNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Stream        classModels.Stream    `gorm:"foreignKey:StreamNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type TimetableController struct {
	gorm.Model
	Number       string `json:"number" gorm:"size:36;unique;not null"`
	SchoolNumber string `json:"school_number" gorm:"size:36;not null"`
	LevelNumber  string `json:"level_number" gorm:"size:36;not null"`
	Type         string `json:"type" gorm:"size:36;not null"`
	Day          string `json:"day" gorm:"size:36;not null"`
	StartTime    string `json:"start_time" gorm:"size:10;not null"`
	EndTime      string `json:"end_time" gorm:"size:10;not null"`
}

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
