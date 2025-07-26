package models

import (
	academicYearModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	subjectModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"time"
)

type (
	SchemeOfWork struct {
		ID                 uint      `gorm:"primaryKey"`
		Number             string    `json:"number" gorm:"size:36;unique;not null"`
		Name               string    `json:"name" form:"name" binding:"required" gorm:"size:200;not null"`
		Status             string    `json:"status" form:"status" binding:"required" gorm:"size:15;not null"`
		CampusClassNumber  string    `json:"campus_class_number" gorm:"size:36;not null"`
		SubjectNumber      string    `json:"subject_number" gorm:"size:36;not null"`
		UserNumber         string    `json:"user_number" gorm:"size:36;not null"`
		AcademicYearNumber string    `json:"academic_year_number" gorm:"size:36;not null"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time
		AcademicYear       academicYearModels.AcademicYear `gorm:"foreignKey:AcademicYearNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		CampusClass        classModels.CampusClass         `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Subject            subjectModel.Subject            `gorm:"foreignKey:SubjectNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User               models.User                     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *SchemeOfWork) TableName() string {
	return "scheme_of_works"
}
