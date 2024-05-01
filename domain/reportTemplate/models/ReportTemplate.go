package models

import (
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type ReportTemplate struct {
	gorm.Model
	Name         string             `gorm:"size:250;not null"`
	Number       string             `gorm:"size:36;unique;not null"`
	SchoolNumber string             `gorm:"size:36;not null"`
	School       schoolModel.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Status       string             `gorm:"size:36;not null"`
	UserNumber   string             `gorm:"size:36;not null"`
	User         userModel.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

func (reportTemplate *ReportTemplate) TableName() string {
	return "report_templates"
}
