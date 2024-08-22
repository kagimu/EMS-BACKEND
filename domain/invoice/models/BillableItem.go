package models

import itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"

type BillableItem struct {
	Number         string          `json:"number" gorm:"size:36;primaryKey;not null"`
	ItemNumber     string          `json:"item_number" gorm:"size:36;not null"`
	Item           itemModels.Item `gorm:"foreignKey:ItemNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Quantity       float64         `json:"quantity" gorm:"not null"`
	UnitPrice      float64         `json:"unit_price" gorm:"not null"`
	BillableNumber string          `json:"billable_number" gorm:"size:36;not null"`
	Billable       Billable        `gorm:"foreignKey:BillableNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

func (billableItem *BillableItem) TableName() string {
	return "billable_items"
}
