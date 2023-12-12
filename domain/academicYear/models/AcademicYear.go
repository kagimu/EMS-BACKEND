package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"gorm.io/gorm"
	"time"
)

type (
	AcademicYear struct {
		gorm.Model
		Number                string                    `json:"number" gorm:"size:36;unique;not null"`
		Code                  int64                     `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		AcademicYearName      string                    `json:"academic_year_name" form:"name" binding:"required" gorm:"size:100;not null"`
		CalendarNumber        string                    `json:"calendar_number" gorm:"size:36;not null"`
		NoOfStudyPeriod       int32                     `json:"no_of_study_period" form:"name" binding:"required" gorm:"not null"`
		AcademicYearStartDate time.Time                 `json:"academic_year_start_date" form:"name" binding:"required" gorm:"not null"`
		AcademicYearEndDate   time.Time                 `json:"academic_year_end_date" form:"name" binding:"required" gorm:"not null"`
		UserNumber            string                    `json:"user_number" gorm:"size:36;not null"`
		CampusLevelNumber     string                    `json:"campus_level_number" gorm:"size:36;not null"`
		User                  models.User               `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		CampusLevel           models2.CampusLevel       `gorm:"foreignKey:CampusLevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		GeneralCalendar       eosModels.GeneralCalendar `gorm:"foreignKey:CalendarNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
