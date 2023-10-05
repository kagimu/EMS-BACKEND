package models

import (
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	subjectModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"gorm.io/gorm"
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
	Campus        campusModels.Campus   `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Subject       subjectModels.Subject `gorm:"foreignKey:SubjectNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Stream        classModels.Stream    `gorm:"foreignKey:StreamNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
