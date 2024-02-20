package models

import "gorm.io/gorm"

type InvoiceItem struct {
	gorm.Model
	ID            uint64 `gorm:"primary_key;auto_increment" json:"id"`
	InvoiceNumber string `json:"invoice_number" gorm:"size:36;not null"`
	ItemNumber    string `json:"item_number" gorm:"size:36;not null"`
	Quantity      int    `json:"quantity" gorm:"not null"`
	UnitPrice     int    `json:"unit_price" gorm:"not null"`
	TotalPrice    int    `json:"total_price" gorm:"-"`
	ItemName      string `json:"item_name" gorm:"-"`
}

// Use after find hook to calculate total price and item name
func (invoiceItem *InvoiceItem) AfterFind(tx *gorm.DB) (err error) {
	invoiceItem.TotalPrice = invoiceItem.Quantity * invoiceItem.UnitPrice
	// find item name from table items

	return
}
