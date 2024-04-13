package models

import (
	"time"

	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	streamModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	subjectModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"gorm.io/gorm"
)

type StudyingTimeTable struct {
	gorm.Model
	Number                  string                `json:"number" gorm:"size:36;not null"`
	StudyingTimeBlockNumber string                `json:"studying_time_block_number" gorm:"size:36;not null"`
	StudyingTimeBlock       StudyingTimeBlock     `gorm:"foreignKey:StudyingTimeBlockNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StreamNumber            string                `json:"stream_number" gorm:"size:36;not null"`
	Day                     time.Time             `json:"day"`
	Stream                  streamModels.Stream   `gorm:"foreignKey:StreamNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SubjectNumber           string                `json:"subject_number" gorm:"size:36;not null"`
	Subject                 subjectModels.Subject `gorm:"foreignKey:SubjectNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StartTime               time.Time             `json:"start_time" gorm:"size:36;"`
	EndTime                 time.Time             `json:"end_time" gorm:"size:36;"`
	UserNumber              string                `json:"user_number" gorm:"size:36;not null"`
	User                    authModel.User        `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
