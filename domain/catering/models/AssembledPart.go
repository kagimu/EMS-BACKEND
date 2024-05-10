package models

type AssembledPart struct {
	FromStockingNumber string   `json:"from_stocking_number" gorm:"size:36;not null"`
	FromStocking       Stocking `gorm:"foreignKey:FromStockingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	NewStockingNumber  string   `json:"new_stocking_number" gorm:"size:36;not null"`
	NewStocking        Stocking `gorm:"foreignKey:NewStockingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

func (assembledPart *AssembledPart) TableName() string {
	return "assembled_parts"
}
