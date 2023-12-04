package models

import (
	academicYearModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"gorm.io/gorm"
	"time"
)

type (
	StudyPeriod struct {
		gorm.Model
		Number               string                          `json:"number" gorm:"size:36;unique;not null"`
		Code                 int64                           `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		StudyPeriodName      string                          `json:"study_period_name" form:"name" binding:"required" gorm:"size:100;not null"`
		StudyPeriodStartDate time.Time                       `json:"academic_year_start_date" form:"name" binding:"required" gorm:"not null"`
		StudyPeriodEndDate   time.Time                       `json:"academic_year_end_date" form:"name" binding:"required" gorm:"not null"`
		AcademicYearNumber   string                          `json:"academic_year_number" gorm:"size:36;not null"`
		CalendarEventNumber  string                          `json:"calendar_event_number" gorm:"size:36;not null"`
		UserNumber           string                          `json:"user_number" gorm:"size:36;not null"`
		LevelNumber          string                          `json:"level_number" gorm:"size:36;not null"`
		User                 models.User                     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		AcademicYear         academicYearModels.AcademicYear `gorm:"foreignKey:AcademicYearNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Level                models2.Level                   `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		CalendarEvent        eosModels.CalendarEvent         `gorm:"foreignKey:CalendarEventNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
