package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	studyPeriod "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

type Enrollment struct {
	gorm.Model
	Number            string                  `json:"number" gorm:"size:36;unique;not null"`
	StudentNumber     string                  `json:"student_number" gorm:"size:36;not null"`
	Student           models.Student          `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriodNumber string                  `json:"study_period_number" gorm:"size:36;not null;uniqueIndex"`
	StudyPeriod       studyPeriod.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber      string                  `json:"campus_number" gorm:"size:36;not null;uniqueIndex"`
	Campus            campusModel.Campus      `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SchoolNumber      string                  `json:"school_number"  gorm:"size:36;not null;uniqueIndex"`
	School            schoolModel.School      `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	EnrollmentDate    string                  `json:"enrollment_date" gorm:"size:15;"`
	Status            string                  `json:"status" gorm:"size:20;"`
	UserNumber        string                  `json:"user_number" gorm:"size:36;index;not null"`
	User              authModels.User         `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
