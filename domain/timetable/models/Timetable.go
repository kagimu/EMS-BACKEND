package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

type TimetableSlots struct {
	gorm.Model
	Name         string `json:"name" gorm:"size:100;not null"`
	Number       string `json:"number" gorm:"size:36;unique;not null"`
	SchoolNumber string `json:"school_number" gorm:"size:36;not null"`
	LevelNumber  string `json:"level_number" gorm:"size:36;not null"`
	Type         string `json:"type" gorm:"size:36;not null"`
	StartTime    string `json:"start_time" gorm:"size:10;not null"`
	EndTime      string `json:"end_time" gorm:"size:10;not null"`
}

type TimetableTask struct {
	gorm.Model
	Number       string `json:"number" gorm:"size:36;unique;not null"`
	Name         string `json:"name" gorm:"size:100;not null"`
	Type         string `json:"type" gorm:"size:36;not null"`
	UserNumber   string `json:"user_number" gorm:"size:36;not null"`
	ClassNumber  string `json:"class_number" gorm:"size:36;not null"`
	StreamNumber string `json:"stream_number" gorm:"size:36;not null"`
}

type StudyTimetable struct {
	gorm.Model
	Number            string             `json:"number" gorm:"size:36;unique;not null"`
	Name              string             `json:"name"`
	StudyPeriodNumber string             `json:"study_period_number" gorm:"size:36;not null"`
	Status            string             `json:"status" gorm:"size:36;not null"`
	UserNumber        string             `json:"user_number" gorm:"size:36;not null"`
	SchoolNumber      string             `json:"school_number" gorm:"size:36;not null"`
	LevelNumber       string             `json:"level_number" gorm:"size:36;not null"`
	CampusNumber      string             `json:"campus_number" gorm:"size:36;not null"`
	ClassNumber       string             `json:"class_number" gorm:"size:36;not null"`
	StreamNumber      string             `json:"stream_number" gorm:"size:36;not null"`
	Stream            classModels.Stream `json:"stream" gorm:"foreignkey:StreamNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriod       models.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type StudyTimetablePeriod struct {
	gorm.Model
	StudyTimetableNumber string             `json:"study_timetable_number" gorm:"size:36;not null"`
	Number               string             `json:"number" gorm:"size:36;unique;not null"`
	StreamNumber         string             `json:"stream_number" gorm:"size:36;not null"`
	Day                  string             `json:"day" gorm:"size:36;not null"`
	DayCode              int32              `json:"day_code" gorm:"size:36;not null"`
	StartTime            string             `json:"start_time" gorm:"size:10;not null"`
	EndTime              string             `json:"end_time" gorm:"size:10;not null"`
	TaskType             string             `json:"task_type" gorm:"size:100;not null"`
	UserNumber           string             `json:"user_number" gorm:"size:36;not null"`
	User                 authModels.User    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Stream               classModels.Stream `gorm:"foreignKey:StreamNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyTimetable       StudyTimetable     `gorm:"foreignKey:StudyTimetableNumber;references:Number; constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
