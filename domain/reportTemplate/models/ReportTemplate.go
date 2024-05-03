package models

import (
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type ReportTemplate struct {
	gorm.Model
	Number       string             `json:"number" gorm:"size:36;unique;not null"`
	Name         string             `json:"name" gorm:"size:250;not null"`
	SchoolNumber string             `json:"school_number" gorm:"size:36;not null"`
	School       schoolModel.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Status       string             `json:"status" gorm:"size:36;not null"`
	UserNumber   string             `json:"user_number" gorm:"size:36;not null"`
	User         userModel.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

func (reportTemplate *ReportTemplate) TableName() string {
	return "report_templates"
}
