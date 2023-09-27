package models

import (
	classModels "github.com/hamsoftug/wholesomeEd/domain/class/models"
	subjectModels "github.com/hamsoftug/wholesomeEd/domain/subject/models"
	"gorm.io/gorm"
	"time"
)

type Timetable struct {
	gorm.Model
	Number        string                `json:"number" gorm:"size:100;unique;not null"`
	StreamNumber  string                `json:"stream_number" gorm:"size:40;not null"`
	Day           string                `json:"day" gorm:"size:40;not null"`
	StartTime     time.Time             `json:"start_time" gorm:"type:time;not null"`
	EndTime       time.Time             `json:"end_time" gorm:"type:time;not null"`
	SubjectNumber string                `json:"subject_number" gorm:"size:40;not null"`
	Subject       subjectModels.Subject `gorm:"foreignKey:SubjectNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Stream        classModels.Stream    `gorm:"foreignKey:StreamNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
