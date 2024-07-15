package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type MenuItem struct {
	gorm.Model
	Number                      string          `json:"number" gorm:"size:36;unique;not null"`
	Name                        string          `json:"name" gorm:"size:50;"`
	Code                        int64           `json:"code"`
	BaseUnitCost                float64         `json:"base_unit_cost"`
	PurchaseUnitOfMeasureNumber string          `json:"purchase_unit_of_measure_number" gorm:"size:36;"`
	PurchaseUnitOfMeasure       UnitOfMeasure   `gorm:"foreignKey:PurchaseUnitOfMeasureNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ConversionToPurchaseUnit    float64         `json:"conversion_to_purchase_unit"`
	Type                        string          `json:"type" gorm:"size:50;not null"`
	SchoolNumber                string          `json:"school_number" gorm:"size:36;not null"`
	School                      models.School   `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UnitOfMeasureNumber         string          `json:"unit_of_measure_number" gorm:"size:36;"`
	UnitOfMeasure               UnitOfMeasure   `gorm:"foreignKey:UnitOfMeasureNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber                  string          `json:"user_number" gorm:"size:36;index;not null"`
	User                        authModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Status                      string          `json:"status" gorm:"size:100;not null"`
}
