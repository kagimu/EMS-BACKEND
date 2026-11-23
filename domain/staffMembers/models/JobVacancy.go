package models

import (
	"time"

	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	departmentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/departments/models"
	"gorm.io/gorm"
)

type JobVacancy struct {
	gorm.Model
	Number            string                      `json:"number" gorm:"size:36;unique;index;not null"`
	JobPositionNumber string                      `json:"job_position_number" gorm:"size:36;not null"`
	JobPosition       JobPosition                 `gorm:"foreignKey:JobPositionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	VacancyTitle      string                      `json:"vacancy_title" gorm:"size:150;not null"`
	ReferenceNumber   string                      `json:"reference_number" gorm:"size:100"`
	DepartmentNumber  string                      `json:"department_number" gorm:"size:36;not null"`
	Department        departmentModels.Department `gorm:"foreignKey:DepartmentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SchoolNumber      string                      `json:"school_number" gorm:"size:36;not null"`
	School            schoolModels.School         `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	BranchLocation    string                      `json:"branch_location" gorm:"size:36"`
	Campus            campusModels.Campus         `gorm:"foreignKey:BranchLocation;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	NumberOfPositions int                         `json:"number_of_positions" gorm:"default:1"`
	AppointmentTerms  string                      `json:"appointment_terms" gorm:"size:100"`
	OpeningDate       time.Time                   `json:"opening_date" gorm:"type:date"`
	ClosingDate       time.Time                   `json:"closing_date" gorm:"type:date"`
	Status            string                      `json:"status" gorm:"size:50;default:'Open'"`
	UserNumber        string                      `json:"user_number" gorm:"size:36;not null"`
	User              authModels.User             `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
