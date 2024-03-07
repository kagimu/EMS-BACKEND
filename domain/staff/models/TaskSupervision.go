package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskSupervision struct {
	gorm.Model
	Number              string    `json:"number" gorm:"size:36;unique;not null"`
	DateOfSupervision   time.Time `json:"date_of_supervision"`
	DescriptiveAnalysis string    `json:"descriptive_analysis" gorm:"type:text"`
	NameOfSupervisor    string    `json:"name_of_supervisor" gorm:"size:30"`
	TaskNumber          string    `json:"task_number" gorm:"size:36"`
	Task                Task      `gorm:"foreignKey:TaskNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
