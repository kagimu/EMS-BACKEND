package models

import (
	academicYearModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"gorm.io/gorm"
	"time"
)

type (
	StudyPeriod struct {
		gorm.Model
		Number                string                          `json:"number" gorm:"size:36;unique;not null"`
		Code                  int64                           `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		StudyPeriodNameNumber string                          `json:"study_period_name_number" gorm:"size:36;not null"`
		StudyPeriodName       StudyPeriodName                 `gorm:"foreignKey:StudyPeriodNameNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
		StudyPeriodStartDate  time.Time                       `gorm:"not null;column:study_period_start_date"`
		StudyPeriodEndDate    time.Time                       `gorm:"not null;column:study_period_end_date"`
		AcademicYearNumber    string                          `gorm:"size:36;not null;column:academic_year_number"`
		CalendarEventNumber   string                          `gorm:"size:36;not null;unique;column:calendar_event_number"`
		UserNumber            string                          `json:"user_number" gorm:"size:36;not null"`
		User                  models.User                     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
		AcademicYear          academicYearModels.AcademicYear `gorm:"foreignKey:AcademicYearNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
		CalendarEvent         eosModels.CalendarEvent         `gorm:"foreignKey:CalendarEventNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	}
)
