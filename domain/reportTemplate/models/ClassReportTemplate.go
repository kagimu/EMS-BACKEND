package models

import (
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusClassModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
)

type ClassReportTemplate struct {
	CampusClassNumber    string                       `gorm:"size:36;not null"`
	CampusClass          campusClassModel.CampusClass `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ReportTemplateNumber string                       `gorm:"size:36;not null"`
	ReportTemplate       ReportTemplate               `gorm:"foreignKey:ReportTemplateNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Number               string                       `gorm:"size:36;unique;not null"`
	UserNumber           string                       `gorm:"size:36;not null"`
	User                 userModel.User               `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
