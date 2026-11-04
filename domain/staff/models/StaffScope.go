package models

import (
	"time"
)

type StaffScope struct {
	ID                   uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Number               string    `json:"number" gorm:"size:36;unique;not null"`
	StaffNumber          string    `json:"staff_number" gorm:"size:36;index;not null"`
	UserNumber           string    `json:"user_number" gorm:"size:36;index"`
	SchoolNumber         string    `json:"school_number" gorm:"size:36;index;not null"`
	CampusNumber         string    `json:"campus_number" gorm:"size:36;index"`
	CurriculumNumber     string    `json:"curriculum_number" gorm:"size:36;index"`
	CampusLevelNumber    string    `json:"campus_level_number" gorm:"size:36;index"`
	ClassNumber          string    `json:"class_number" gorm:"size:36;index"`
	StreamNumber         string    `json:"stream_number" gorm:"size:36;index"`
	AssessmentAreaNumber string    `json:"assessment_area_number" gorm:"size:36;index"`
	AddedBy              string    `json:"added_by" gorm:"size:36"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (s *StaffScope) TableName() string {
	return "staff_scopes"
}
