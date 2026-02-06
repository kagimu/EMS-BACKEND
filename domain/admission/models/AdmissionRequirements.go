package models

import (
	"time"

	"gorm.io/gorm"
)

type AdmissionRequirement struct {
	gorm.Model
	Number             string    `json:"number" gorm:"size:36;primaryKey"`
	ItemNumber         string    `json:"item_number" gorm:"size:36;not null"`
	CampusClassNumber  string    `json:"campus_class_number" gorm:"size:36;"`
	StudyPeriodNumber  string    `json:"study_period_number" gorm:"size:36;"`
	AcademicYearNumber string    `json:"academic_year_number" gorm:"size:36;not null"`
	UserNumber         string    `json:"user_number" gorm:"size:36;not null"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	SchoolNumber       string    `json:"school_number" gorm:"size:36;"`
}
