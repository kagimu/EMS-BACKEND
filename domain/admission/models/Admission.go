package models

import (
	academicYearModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

type Admission struct {
	gorm.Model
	Number             string                         `json:"number" gorm:"size:36;unique;not null"`
	StudentNumber      string                         `json:"student_number" gorm:"size:36;unique;not null"`
	Student            studentModel.Student           `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber       string                         `json:"campus_number" gorm:"size:36;"`
	Campus             campusModels.Campus            `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusClassNumber  string                         `json:"campus_class_number" gorm:"size:36;"`
	CampusClass        classModels.CampusClass        `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriodNumber  string                         `json:"study_period_number" gorm:"size:36;"`
	StudyPeriod        studyPeriodModel.StudyPeriod   `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AcademicYearNumber string                         `json:"academic_year_number" gorm:"size:36;not null"`
	AcademicYear       academicYearModel.AcademicYear `gorm:"foreignKey:AcademicYearNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StreamNumber       string                         `json:"stream_number" gorm:"size:36;"`
	Stream             classModels.Stream             `gorm:"foreignKey:StreamNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber         string                         `json:"user_number" gorm:"size:36;not null"`
	User               authModels.User                `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Status             string                         `json:"status" gorm:"size:20;default:'PENDING'"`
	SchoolNumber       string                         `json:"school_number" gorm:"size:36;"`
	School             schoolModels.School            `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
