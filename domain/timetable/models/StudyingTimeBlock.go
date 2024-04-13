package models

import (
	"time"

	academicYearModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	levelModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"gorm.io/gorm"
)

type StudyingTimeBlock struct {
	gorm.Model
	Number             string                         `json:"number" gorm:"size:36;not null"`
	StartTime          time.Time                      `json:"start_time"`
	EndTime            time.Time                      `json:"end_time"`
	Length             time.Duration                  `json:"length"`
	AcademicYearNumber string                         `json:"academic_year_number" gorm:"size:36;"`
	AcademicYear       academicYearModel.AcademicYear `gorm:"foreignKey:AcademicYearNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber       string                         `json:"campus_number" gorm:"size:36;"`
	Campus             campusModel.Campus             `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	LevelNumber        string                         `json:"level_number" gorm:"size:36;"`
	Level              levelModel.Level               `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber         string                         `json:"user_number" gorm:"size:36;not null"`
	User               authModel.User                 `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
