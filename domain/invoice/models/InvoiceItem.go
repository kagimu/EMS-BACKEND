package models

import (
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
)

type InvoiceItem struct {
	ID            uint            `gorm:"primarykey"`
	Number        string          `json:"number" gorm:"size:36;unique;not null"`
	InvoiceNumber string          `json:"invoice_number" gorm:"size:36;not null"`
	Invoice       Invoice         `gorm:"foreignKey:InvoiceNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ItemNumber    string          `json:"item_number" gorm:"size:36;not null"`
	Item          itemModels.Item `gorm:"foreignKey:ItemNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Quantity      float64         `json:"quantity" gorm:"not null"`
	UnitPrice     float64         `json:"unit_price" gorm:"not null"`
}
