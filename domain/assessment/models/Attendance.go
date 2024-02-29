package models

import (
	"time"

	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	Number            string                       `json:"number" gorm:"size:36;unique;not null"`
	DayOfTheWeek      time.Time                    `json:"day_of_the_week" gorm:"size:36;"`
	CheckIn           time.Time                    `json:"check_in"`
	CheckOut          time.Time                    `json:"check_out"`
	Present           bool                         `json:"present"`
	ClassNumber       string                       `json:"class_number" gorm:"size:36;"`
	Class             classModel.Class             `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SchoolNumber      string                       `json:"school_number" gorm:"size:36;"`
	School            schoolModel.School           `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber      string                       `json:"campus_number" gorm:"size:36;"`
	Campus            campusModel.Campus           `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriodNumber string                       `json:"study_period_number" gorm:"size:36;"`
	StudyPeriod       studyPeriodModel.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudentNumber     string                       `json:"student_number" gorm:"size:36;index;"`
	Student           studentModel.Student         `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber        string                       `json:"user_number" gorm:"size:36;index;not null"`
	User              authModel.User               `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
