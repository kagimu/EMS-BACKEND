package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type ItemAssembly struct {
	gorm.Model
	Number              string          `json:"number" gorm:"size:36;unique;not null"`
	Name                string          `json:"name" gorm:"size:50;"`
	Code                uint64          `json:"code" gorm:"unique;not null"`
	Type                string          `json:"type" gorm:"size:50;not null"`
	CampusNumber        string          `json:"campus_number" gorm:"size:36;not null"`
	UnitOfMeasureNumber string          `json:"unit_of_measure_number" gorm:"size:36;not null"`
	PrimaryItemNumber   string          `json:"primary_item_number" gorm:"size:36;not null"`
	PrimaryItem         MenuItem        `gorm:"foreignKey:PrimaryItemNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SecondaryItemNumber string          `json:"secondary_item_number" gorm:"size:36;not null"`
	SecondaryItem       MenuItem        `gorm:"foreignKey:SecondaryItemNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Quantity            float64         `json:"quantity" gorm:"not null"`
	UserNumber          string          `json:"user_number" gorm:"size:36;index;not null"`
	User                authModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
