package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"gorm.io/gorm"
)

type MenuItem struct {
	gorm.Model
	Number              string              `json:"number" gorm:"size:36;unique;not null"`
	Name                string              `json:"name" gorm:"size:50;"`
	Code                uint64              `json:"code"`
	Type                string              `json:"type"`
	CampusNumber        string              `json:"campus_number" gorm:"size:36;"`
	Campus              campusModels.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UnitOfMeasureNumber string              `json:"unit_of_measure_number" gorm:"size:36;"`
	UnitOfMeasure       UnitOfMeasure       `gorm:"foreignKey:UnitOfMeasureNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber          string              `json:"user_number" gorm:"size:36;index;not null"`
	User                authModels.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ItemAssemblies      []ItemAssembly      `gorm:"foreignKey:PrimaryItem;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
