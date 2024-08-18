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
		StudyPeriodStartDate  time.Time                       `json:"academic_year_start_date" form:"name" binding:"required" gorm:"not null"`
		StudyPeriodEndDate    time.Time                       `json:"academic_year_end_date" form:"name" binding:"required" gorm:"not null"`
		AcademicYearNumber    string                          `json:"academic_year_number" gorm:"size:36;not null"`
		CalendarEventNumber   string                          `json:"calendar_event_number" gorm:"size:36;not null;unique"`
		UserNumber            string                          `json:"user_number" gorm:"size:36;not null"`
		User                  models.User                     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
		AcademicYear          academicYearModels.AcademicYear `gorm:"foreignKey:AcademicYearNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
		CalendarEvent         eosModels.CalendarEvent         `gorm:"foreignKey:CalendarEventNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	}
)
