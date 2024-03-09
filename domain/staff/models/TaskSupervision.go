package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type TaskSupervision struct {
	gorm.Model
	Number              string      `json:"number" gorm:"size:36;unique;not null"`
	DateOfSupervision   string      `json:"date_of_supervision"`
	DescriptiveAnalysis string      `json:"descriptive_analysis" gorm:"type:text"`
	NameOfSupervisor    string      `json:"name_of_supervisor" gorm:"size:30"`
	TaskNumber          string      `json:"task_number" gorm:"size:36"`
	Task                Task        `gorm:"foreignKey:TaskNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber          string      `json:"user_number" gorm:"size:36;index;not null"`
	User                models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
